import { ssh, scpFrom } from './helpers/ssh'
import { env } from './helpers/env'
import * as path from 'node:path'
import * as fs from 'node:fs'
import { execSync } from 'node:child_process'

const TMP_DIR = '/tmp/syncloud/peertube-ui'

export default async function () {
  const out = path.join(env('PLAYWRIGHT_ARTIFACT_DIR'), 'playwright', env('PLAYWRIGHT_PROJECT'))
  fs.mkdirSync(out, { recursive: true })

  ssh(`mkdir -p ${TMP_DIR}`, { throw: false })
  ssh(`journalctl > ${TMP_DIR}/journalctl.log`, { throw: false })
  ssh(`cat /var/snap/platform/current/config/authelia/config.yml > ${TMP_DIR}/authelia.config.yml.log`, { throw: false })
  scpFrom(`${TMP_DIR}/*`, out, { throw: false })
  try { execSync(`chmod -R a+r ${out}`) } catch {}
}
