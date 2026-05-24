-- Drop retry tracking from test_results
ALTER TABLE test_results DROP COLUMN IF EXISTS retry_count;

-- Drop schedule_runs
DROP TABLE IF EXISTS schedule_runs;

-- Drop webhooks
DROP TABLE IF EXISTS webhooks;

-- Drop dataset link from test_cases
ALTER TABLE test_cases DROP COLUMN IF EXISTS dataset_id;

-- Drop datasets
DROP TABLE IF EXISTS datasets;

-- Drop test_cases additions
ALTER TABLE test_cases DROP COLUMN IF EXISTS max_retries;
ALTER TABLE test_cases DROP COLUMN IF EXISTS timeout_ms;
ALTER TABLE test_cases DROP COLUMN IF EXISTS pre_script;
ALTER TABLE test_cases DROP COLUMN IF EXISTS post_script;

-- Drop test_suites additions
ALTER TABLE test_suites DROP COLUMN IF EXISTS max_retries;
ALTER TABLE test_suites DROP COLUMN IF EXISTS retry_delay_ms;
ALTER TABLE test_suites DROP COLUMN IF EXISTS parallel;
ALTER TABLE test_suites DROP COLUMN IF EXISTS max_parallel;
ALTER TABLE test_suites DROP COLUMN IF EXISTS timeout_ms;

-- Drop environment_id from test_runs
ALTER TABLE test_runs DROP COLUMN IF EXISTS environment_id;

-- Drop environments
DROP INDEX IF EXISTS idx_environments_project;
DROP TABLE IF EXISTS environments;
