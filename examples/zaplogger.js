import zaplogger from 'k6/x/zaplogger';
import exec from 'k6/execution';
const mylogger = zaplogger.initLogger("./test.log")
export const options = {
  setupTimeout: '30m',
  discardResponseBodies: false,
  scenarios: {
    contacts: {
        executor: 'per-vu-iterations',
        vus: 5,
        iterations: 20,
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
}
export function teardown() {
  mylogger.sync()
}