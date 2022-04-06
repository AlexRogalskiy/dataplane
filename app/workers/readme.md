
### Using NATS CLI

```shell
nats context add nats --server localhost:14222 --description "NATS Dataplane" --select

nats sub <channel>

nats pub <channel> "message"

nats sub js.out.testing --ack 

# wildcard use
nats sub "workertask.*"
```

### Measure task execution time in milliseconds in Postgresql
```sql
SELECT
  task_id,
  "status",
  start_dt,
  end_dt,
  Extract(epoch FROM (end_dt - start_dt))*1000 AS milliseconds,
  Extract(epoch FROM (end_dt - start_dt)) AS seconds
FROM worker_tasks 
where end_dt >'0001-01-01'
order by created_at desc;
```

### Measure pipeline execution time
```sql
SELECT
  run_id,
  "status",
  created_at,
  ended_at,
  Extract(epoch FROM (ended_at - created_at))*1000 AS milliseconds,
  Extract(epoch FROM (ended_at - created_at)) AS seconds
FROM pipeline_runs 
where ended_at >'0001-01-01'
order by created_at desc;
```

### Long running bash test for parallel runs
```shell
bash:
for((i=1;i<=10; i+=1)); do echo "2nd run $i times"; sleep 0.5; done

sh:
i=0; while [ $(($i)) -le 25 ]; do i=$(($i + 1)); sleep 0.5; echo "run $i times"; done;
```

```python
import time
for x in range(0, 10):
    print("We're on %d" % (x))
    time.sleep(2)
```

### Number of database connections
```sql
SELECT * FROM pg_stat_database;
SELECT sum(numbackends) FROM pg_stat_database;

select * from pg_stat_activity;
```


