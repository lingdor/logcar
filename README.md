# logcar
a log middleware to direct streams, rolling, filter, split to files  kafka, hdfs or others.

case 1:
```text
|--progress---|                |-----------------|
|  app        | ---pipline---> |    logcar       |
|-------------|                |-----------------|
```
case 2:
```text

           Dev  /  Ops
|--container--|    |----------|     |-----------------|    |----------------------|
|  app        | -> |  logcar  | ->  |  kafka          | -> | logcar file-appender |
|(go|java|...)|    |----------|     |-----------------|    |----------------------|
|-------------|                            |               |----------------------|
                                           |-------------> |        flink/ES      |
                                                           |----------------------|



```
# Command

```shell
./xx-server|logcar -f ./logcar.yml
```




