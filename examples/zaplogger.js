import zaplogger from 'k6/x/zaplogger';
import exec from 'k6/execution';
// const mylogger = zaplogger.initLogger("./test.log")
// mylogger.infow("msg", "key1", `start`)
export function setup() {
}
export const options = {
    setupTimeout: '30m',
    discardResponseBodies: false,
    scenarios: {
      contacts: {
          executor: 'per-vu-iterations',
          vus: 2,
          iterations: 10,
          maxDuration: '5m',
          exec: 'test2',
          tags: { my_custom_tag: 'mytag' },
          env: { MYVAR: 'contacts' },
        }
    },
  };
export function test2(){
    const mylogger = zaplogger.initLogger("./test" +`${exec.vu.idInTest}` +".log")
    mylogger.infow("msg", "key1", `${exec.vu.idInTest}`)
}
export function teardown(data) {
}