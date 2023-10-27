# xk6-zap

This is a [k6](https://github.com/grafana/k6) extension using the
[xk6](https://github.com/grafana/xk6) system.

Support logging to file.


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
  xk6 build --with github.com/Mistyrain520/xk6-zap@latest
  ```




## Example

```javascript
// script.js
import zaplogger from 'k6/x/zaplogger';
import { sleep } from 'k6';

const mylogger = zaplogger.initLogger("./test.log")

export default function () { 
  mylogger.infow("msg", "key", "gagga")
  mylogger.infow("key", "gagga", "key1", "values1")
  sleep(5)
  }
export function teardown() {
  mylogger.sync()
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


## TUDO
- [ ] 文件不存在的时候，并发写入有问题；文件如果已经存在了，就没问题。奇怪
