import zaplogger from 'k6/x/zaplogger';
import { sleep } from 'k6';

const mylogger = zaplogger.initLogger("./test.log")
console.log("@@@@@@@@@@")
export function setup() {
  // mylogger.sync()
}
export default function () { 
  const aa = zaplogger.zapObject("message", "mkey", "mvalue", "mkey2", 3)
  // mylogger.infow("msg", "key", "gagga")
  mylogger.infow("msg", "key1", "values1", aa)
  sleep(5)
  }
export function teardown() {
  // mylogger.sync()
}