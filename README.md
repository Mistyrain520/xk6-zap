# xk6-zap

This is a [k6](https://github.com/grafana/k6) extension using the
[xk6](https://github.com/grafana/xk6) system.

Support logging to file.

## v2.0.2

- Fixed concurrent writes when multiple VUs create the same log file.
- Fixed concurrent log rotation causing split logs or Windows rename conflicts.
- Added `closeLogger()` and `closeAll()` for releasing log file handles.


## Build

To build a `k6` binary with this plugin, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- [xk6](https://github.com/grafana/xk6)
- Git

Then:

1. Install `xk6`:
  ```shell
  go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  xk6 build v2.0.0 --with github.com/Mistyrain520/xk6-zap/v2@v2.0.2
  ```




## Example

```javascript
// script.js
import zaplogger from 'k6/x/zaplogger';
import { sleep } from 'k6';
// Multiple VUs can safely share the same log file.
const mylogger = zaplogger.initLogger("./test.log")
mylogger.infow("msg", "key1", `start`)
export default function () { 
  mylogger.infow("msg", "key", "gagga")
  mylogger.infow("key", "gagga", "key1", "values1")
  sleep(5)
  }
// Different VUs can also use different files.
// export default function () { 
//  const mylogger = zaplogger.initLogger("./test" +`${exec.vu.idInTest}` +".log")
//   mylogger.infow("msg", "key", "gagga")
//   mylogger.infow("key", "gagga", "key1", "values1")
//   sleep(5)
//   }
export function teardown() {
  mylogger.sync()
  zaplogger.closeAll()
}
```

Result output:

```shell
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: .\examples\zaplogger.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)


     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=15.65ms min=15.65ms med=15.65ms max=15.65ms p(90)=15.65ms p(95)=15.65ms
     iterations...........: 1   63.864709/s


running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU
```
Then you can see that a file is generated locally
```
//test.log
{"level":"INFO","ts":"2023-10-26T23:44:21.858+0800","msg":"msg","key":"gagga"}
{"level":"INFO","ts":"2023-10-26T23:44:21.874+0800","msg":"msg","key":"gagga","key1":"values1"}

```


## Fixed: concurrent logging and rotation

`initLogger()` shares one managed writer for each canonical file path.
Concurrent VUs can safely create the same file and rotate it without opening
competing file handles or racing on backup names.

Call `zaplogger.closeAll()` in `teardown()` to release all file handles:

```javascript
export function teardown() {
  zaplogger.closeAll()
}
```

## Verify concurrency and rotation

Build a k6 binary containing only this local extension:

```powershell
xk6 build v2.0.0 `
  --output .\k6-zap.exe `
  --with 'github.com/Mistyrain520/xk6-zap/v2=D:\github\xk6-zap'
```

Run the first-creation or existing-file test:

```powershell
$env:VUS = '256'
$env:TARGET = 'D:\github\xk6-zap\shared.log'
.\k6-zap.exe run --quiet .\examples\concurrent_shared_file.js
```

Force rotation with a 1 MB limit and a larger payload:

```powershell
$env:VUS = '512'
$env:MAX_SIZE = '1'
$env:MAX_BACKUPS = '20'
$env:MAX_AGE = '1'
$env:PAYLOAD_SIZE = '4096'
$env:TARGET = 'D:\github\xk6-zap\rotating.log'
.\k6-zap.exe run --quiet .\examples\concurrent_shared_file.js
```

Each JSON line contains `"msg":"concurrent-write"` and a unique `worker`.
The total unique workers across the current and rotated files must equal `VUS`.
