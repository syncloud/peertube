import os
from os.path import join
from subprocess import check_output
import json

import pytest
import requests
from requests.packages.urllib3.exceptions import InsecureRequestWarning
from syncloudlib.http import wait_for_rest
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.installer import local_install
import time

TMP_DIR = '/tmp/syncloud'

requests.packages.urllib3.disable_warnings(InsecureRequestWarning)


@pytest.fixture(scope="session")
def module_setup(request, device, app_dir, artifact_dir):
    def module_teardown():
        device.run_ssh('ls -la /var/snap/peertube/current/config > {0}/config.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/peertube/current/config/webui.yaml {0}/webui.yaml.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/peertube/current/config/authelia/config.yml {0}/authelia.config.yml.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/peertube/current/config/authelia/authrequest.conf {0}/authelia.authrequest.conf.log'.format(TMP_DIR), throw=False)
        device.run_ssh('top -bn 1 -w 500 -c > {0}/top.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ps auxfw > {0}/ps.log'.format(TMP_DIR), throw=False)
        device.run_ssh('netstat -nlp > {0}/netstat.log'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl | tail -1000 > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /snap > {0}/snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /snap/peertube > {0}/snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap > {0}/var.snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap/peertube > {0}/var.snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap/peertube/current/ > {0}/var.snap.current.ls.log'.format(TMP_DIR),
                       throw=False)
        device.run_ssh('ls -la /var/snap/peertube/common > {0}/var.snap.common.ls.log'.format(TMP_DIR),
                       throw=False)
        device.run_ssh('ls -la /data > {0}/data.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /data/peertube > {0}/data.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cat /etc/hosts > {0}/hosts.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cat /var/snap/platform/current/config/authelia/config.yml > {0}/authelia.config.log'.format(TMP_DIR), throw=False)

        app_log_dir = join(artifact_dir, 'log')
        os.mkdir(app_log_dir)
        device.scp_from_device('{0}/*'.format(TMP_DIR), app_log_dir)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, device, device_host, app, domain):
    add_host_alias(app, device_host, domain)
    device.run_ssh('date', retries=100)
    device.run_ssh('mkdir {0}'.format(TMP_DIR))


@pytest.mark.flaky(retries=50, delay=10)
def test_activate_device(device):
    response = retry(device.activate_custom)
    assert response.status_code == 200, response.text


def test_ca_cert(device, app_domain):
    device.run_ssh('CURL_CA_BUNDLE=/var/snap/platform/current/syncloud.ca.crt curl -v https://{0} > {1}/ssl.ca.log 2>&1'.format(app_domain, TMP_DIR))


def test_install(app_archive_path, device_host, device_password):
    local_install(device_host, device_password, app_archive_path)


def test_index(app_domain):
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 100)


def test_yt_dlp_bundled(device):
    config = device.run_ssh('cat /var/snap/peertube/current/config/production.yaml')
    assert "name: 'yt-dlp'" in config
    assert 'python_path: null' in config
    link = device.run_ssh('readlink /data/peertube/bin/yt-dlp').strip()
    assert link == '/snap/peertube/current/peertube/yt-dlp/yt-dlp', link
    device.run_ssh('/data/peertube/bin/yt-dlp --version > {0}/yt-dlp.version.log'.format(TMP_DIR))


@pytest.mark.flaky(retries=30, delay=10)
def test_external_auth_registered(app_domain):
    response = requests.get('https://{0}/api/v1/config'.format(app_domain), verify=False)
    assert response.status_code == 200, response.text
    auths = response.json()['plugin']['registeredExternalAuths']
    assert any(a['npmName'] == 'peertube-plugin-auth-openid-connect' for a in auths), auths


def test_services_not_restarting(device):
    for service in ['nginx', 'peertube', 'postgresql', 'redis']:
        restarts = device.run_ssh('systemctl show snap.peertube.{0}.service -p NRestarts --value'.format(service)).strip()
        assert restarts == '0', '{0} restarted {1} times'.format(service, restarts)


def __log_data_dir(device):
    device.run_ssh('ls -la /data')
    device.run_ssh('mount')
    device.run_ssh('ls -la /data/')
    device.run_ssh('ls -la /data/peertube')


def test_storage_change_event(device):
    device.run_ssh('snap run peertube.storage-change > {0}/storage-change.log'.format(TMP_DIR))


def test_access_change(device):
    device.run_ssh('snap run peertube.access-change > {0}/access-change.log'.format(TMP_DIR))


def test_remove(device, app):
    response = device.app_remove(app)
    assert response.status_code == 200, response.text


def test_reinstall(app_archive_path, device_host, device_password, app_domain):
    local_install(device_host, device_password, app_archive_path)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 100)


def test_upgrade(app_archive_path, device_host, device_password, app_domain):
    local_install(device_host, device_password, app_archive_path)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 100)


@pytest.mark.flaky(retries=3, delay=1)
def test_access_change_after_upgrade(device):
    device.run_ssh("snap run peertube.access-change")


def test_backup(device, artifact_dir):
    device.run_ssh("snap run platform.cli backup create peertube")
    response = device.run_ssh("snap run platform.cli backup list")
    open('{0}/cli.backup.list.json'.format(artifact_dir), 'w').write(response)
    print(response)
    backup = json.loads(response)[0]
    device.run_ssh('tar tf {0}/{1}'.format(backup['path'], backup['file']))
    device.run_ssh("snap run platform.cli backup restore {0}".format(backup['file']))


def test_npm(device, artifact_dir):
    device.run_ssh("snap run peertube.npm run")


def test_sql_plugin(device, artifact_dir):
    device.run_ssh("snap run peertube.psql -U peertube -d peertube -c 'select * from plugin'", retries=10)


def retry(method, retries=10):
    attempt = 0
    exception = None
    while attempt < retries:
        try:
            return method()
        except Exception as e:
            exception = e
            print('error (attempt {0}/{1}): {2}'.format(attempt + 1, retries, str(e)))
            time.sleep(5)
        attempt += 1
    raise exception
