import zaplogger from 'k6/x/zaplogger';
import exec from 'k6/execution';
import { sleep } from 'k6';
const mylogger = zaplogger.initLogger("./test.log")
sleep(3)
export const options = {
  setupTimeout: '30m',
  discardResponseBodies: false,
  scenarios: {
    contacts: {
        executor: 'per-vu-iterations',
        vus: 2,
        iterations: 40,
        maxDuration: '5m',
        exec: 'test2',
        tags: { my_custom_tag: 'mytag' },
        env: { MYVAR: 'contacts' },
      }
  },
};
export function test2() {
  mylogger.infow("msg", "key", "gagga")
  mylogger.infow("msg", "key1", `${exec.vu.idInTest}`)
  mylogger.sync()
}
export function teardown() {
  mylogger.sync()
}