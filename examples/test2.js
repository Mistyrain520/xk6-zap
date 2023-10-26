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