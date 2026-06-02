# ingest-runtime

## ADDED Requirements

### Requirement: usage_log SHALL be partitioned by day and pruned to 90 days
A scheduled maintenance job MUST ensure daily `usage_log` partitions exist ahead of time and MUST
drop partitions older than 90 days, enforcing the §9.4 retention without losing in-window rows.

#### Scenario: Partitions are maintained
- **WHEN** the maintenance job runs
- **THEN** today's and the next several days' partitions exist and rows land in the correct dated partition

#### Scenario: Old data is pruned
- **WHEN** a partition's date range is entirely older than 90 days
- **THEN** that partition is detached and dropped, while in-window data is retained
