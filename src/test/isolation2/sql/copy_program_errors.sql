-- ######################################################
-- Test COPY PROGRAM under error conditions
-- ######################################################

1:CREATE TABLE non_zero_exit_in_child(i int);

-- Take stock of pipes before operations, we will repeat this query to ensure
-- that we don't leak fds. We should see the same output after each test case.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info();

--
-- Test copying from/to program, when the child exits with a non-zero exit code. We
-- should print the exit code, along with stderr.
--
1:COPY non_zero_exit_in_child FROM PROGRAM 'cat /file/does/not/exist';
-- Child shell and grandchild PROGRAM should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with shell and PROGRAM should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info() WHERE content = -1;

1:COPY non_zero_exit_in_child TO PROGRAM 'cat /file/does/not/exist';
-- Child shell and grandchild PROGRAM should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with shell and PROGRAM should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info() WHERE content = -1;

--
-- Test copying from/to program, when the child process is terminated by a signal.
-- We should report termination (of either the intermediate shell OR the
-- PROGRAM, depending on the platform)
--
1&:COPY non_zero_exit_in_child FROM PROGRAM 'sleep 8';
SELECT kill_children(pid, /* TERM */ 15) FROM pg_stat_activity
    WHERE query LIKE '%sleep 8%' AND query NOT LIKE '%pg_stat_activity%' AND state = 'active';
--start_ignore
1<:
--end_ignore
-- Child PROGRAM (and any intermediate shell) should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with PROGRAM (and any intermediate shell) should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info() WHERE content = -1;

1&:COPY non_zero_exit_in_child TO PROGRAM 'sleep 8';
SELECT kill_children(pid, /* TERM */ 15) FROM pg_stat_activity
WHERE query LIKE '%sleep 8%' AND query NOT LIKE '%pg_stat_activity%' AND state = 'active';
--start_ignore
1<:
--end_ignore
-- Child PROGRAM (and any intermediate shell) should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with PROGRAM (and any intermediate shell) should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info() WHERE content = -1;

-- Now do SIGPIPE. Since we do SIG_DFL for descendants of the Postgres backend,
-- the Postgres backend should report that the child was terminated.

1&:COPY non_zero_exit_in_child FROM PROGRAM 'sleep 8';
SELECT kill_children(pid, /* PIPE */ 13) FROM pg_stat_activity
WHERE query LIKE '%sleep 8%' AND query NOT LIKE '%pg_stat_activity%' AND state = 'active';
--start_ignore
1<:
--end_ignore
-- Child PROGRAM (and any intermediate shell) should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with PROGRAM (and any intermediate shell) should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info() WHERE content = -1;

1&:COPY non_zero_exit_in_child TO PROGRAM 'sleep 9';
SELECT kill_children(pid, /* PIPE */ 13) FROM pg_stat_activity
WHERE query LIKE '%sleep 9%' AND query NOT LIKE '%pg_stat_activity%' AND state = 'active';
--start_ignore
1<:
--end_ignore
-- Child PROGRAM (and any intermediate shell) should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with shell and PROGRAM should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info() WHERE content = -1;

--
-- Test copying from program, in the face of data loading errors. Check for fd leaks
-- and zombies on every segment. Sample output of a zombie:
-- postgres=# COPY data_load_error FROM PROGRAM 'cat /tmp/data_load_error_fifo<SEGID>' ON SEGMENT;
-- ERROR:  new row for relation "data_load_error" violates check constraint "data_load_error_i_check"  (seg0 127.0.1.1:7002 pid=735190)
-- DETAIL:  Failing row contains (10).
-- CONTEXT:  COPY data_load_error, line 10: "10"
-- postgres=# select content, (get_descendant_process_info(pid)).* from gp_backend_info();
--  content | pg_backend_pid |                              descendant_proc
-- ---------+----------------+----------------------------------------------------------------------------
--        0 |         743075 | psutil.Process(pid=743159, name='sh', status='zombie', started='14:36:37')
-- (1 row)
--
!\retcode mkfifo /tmp/data_load_error_fifo0;
!\retcode mkfifo /tmp/data_load_error_fifo1;
!\retcode mkfifo /tmp/data_load_error_fifo2;
1:CREATE TABLE data_load_error(i int CHECK (i < 10)) DISTRIBUTED REPLICATED;
1&:COPY data_load_error FROM PROGRAM 'cat /tmp/data_load_error_fifo<SEGID>' ON SEGMENT;
-- Violate the check constraint by writing to one of the named pipes.
!\retcode seq 1 10 > /tmp/data_load_error_fifo0;
1<:
-- Child shell and grandchild PROGRAM should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with shell and PROGRAM should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info();
!\retcode rm /tmp/data_load_error_fifo*;

--
-- Test copying from program, in the face of data loading errors with more data.
-- This is different from the previous case, as with enough data, the child
-- shell and PROGRAM would keep running, even though the COPY errored out. Our
-- code now ensures that such a situation doesn't arise. Example bad output:
-- postgres=# COPY data_load_error FROM PROGRAM 'cat /tmp/data_load_error_fifo<SEGID>' ON SEGMENT;
-- ERROR:  new row for relation "data_load_error" violates check constraint "data_load_error_i_check"  (seg0 127.0.1.1:7002 pid=742816)
-- DETAIL:  Failing row contains (100000).

-- postgres=# select content, (get_descendant_process_info(pid)).* from gp_backend_info();
--  content | pg_backend_pid |                                descendant_proc
-- ---------+----------------+-------------------------------------------------------------------------------
--        0 |         742816 | psutil.Process(pid=742824, name='sh', status='sleeping', started='14:34:23')
--        0 |         742816 | psutil.Process(pid=742825, name='cat', status='sleeping', started='14:34:23')
-- (2 rows)
-- Check fd leaks and zombies on every segment.
--
!\retcode mkfifo /tmp/data_load_error_fifo0;
!\retcode mkfifo /tmp/data_load_error_fifo1;
!\retcode mkfifo /tmp/data_load_error_fifo2;
1:CREATE TABLE data_load_error1(i int CHECK (i < 100000)) DISTRIBUTED REPLICATED;
1&:COPY data_load_error1 FROM PROGRAM 'cat /tmp/data_load_error_fifo<SEGID>' ON SEGMENT;
-- Violate the check constraint by writing to one of the named pipes.
!\retcode seq 1 100000 > /tmp/data_load_error_fifo0;
1<:
-- Child shell and grandchild PROGRAM should be gone
1:SELECT content, (get_descendant_process_info(pid)).* from gp_backend_info();
-- Pipe fds linking parent backend with shell and PROGRAM should be gone.
1:SELECT content, count_child_pipes(pid, content) FROM gp_backend_info();
!\retcode rm /tmp/data_load_error_fifo*;
