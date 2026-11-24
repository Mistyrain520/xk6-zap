import zaplogger from 'k6/x/zaplogger';
import exec from 'k6/execution';
import { sleep } from 'k6';

// Run with TARGET pointing to one shared file. Count the records containing
// "concurrent-write" after the run and compare them with VUS.
const target = __ENV.TARGET || './shared-concurrency.log';
const vus = Number(__ENV.VUS || 256);
const maxSize = Number(__ENV.MAX_SIZE || 500);
const maxBackups = Number(__ENV.MAX_BACKUPS || 3);
const maxAge = Number(__ENV.MAX_AGE || 28);
const payloadSize = Number(__ENV.PAYLOAD_SIZE || 0);

export const options = {
  scenarios: {
    shared_file: {
      executor: 'per-vu-iterations',
      vus,
      iterations: 1,
      maxDuration: '30s',
    },
  },
};

export default function () {
  const logger = zaplogger.initLogger(target, false, maxSize, maxBackups, maxAge);

  // Let every VU finish initialization before all writers open the same file.
  sleep(0.5);
  const fields = ['worker', exec.vu.idInTest];
  if (payloadSize > 0) {
    fields.push('payload', 'x'.repeat(payloadSize));
  }
  logger.infow('concurrent-write', ...fields);
}

export function teardown() {
  zaplogger.closeAll();
}
