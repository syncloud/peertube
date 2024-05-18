import pytest
from os.path import dirname, join
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from subprocess import check_output
from syncloudlib.integration.hosts import add_host_alias
from selenium.webdriver.support import expected_conditions as EC

DIR = dirname(__file__)
TMP_DIR = '/tmp/syncloud/ui'


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir, ui_mode):
    def module_teardown():
        device.activated()
        device.run_ssh('mkdir -p {0}'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl > {0}/journalctl.ui.{1}.log'.format(TMP_DIR, ui_mode), throw=False)
        device.run_ssh('cat /var/snap/platform/current/config/authelia/config.yml > {0}/authelia.config.ui.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), join(artifact_dir, 'log'))
        check_output('cp /videos/* {0}'.format(artifact_dir), shell=True)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, app, domain, device_host):
    add_host_alias(app, device_host, domain)


def test_login(selenium, device_user, device_password):
    selenium.driver.get("https://auth.buster.com")
    selenium.find_by(By.ID, "username-textfield")
    selenium.screenshot('auth')
    selenium.open_app()
    selenium.find_by(By.XPATH, "//a[contains(.,'Login')]").click()
    selenium.find_by(By.XPATH, "//a[contains(.,'My Syncloud')]").click()
    selenium.find_by(By.ID, "username-textfield").send_keys(device_user)
    password = selenium.find_by(By.ID, "password-textfield")
    password.send_keys(device_password)
    selenium.screenshot('login')
    #password.send_keys(Keys.RETURN)
    selenium.find_by(By.ID, "sign-in-button").click()
    selenium.find_by(By.ID, "accept-button").click()
    selenium.find_by(By.CLASS_NAME, "publish-button-label")
    selenium.screenshot('main')

def test_publish_video(selenium):
    selenium.find_by(By.XPATH, "//input[@type='button' and @value='Close']").click()
    selenium.find_by(By.CLASS_NAME, "publish-button-label").click()
    #selenium.find_by_xpath("//span[text()='New post']").click()
    #selenium.find_by_xpath("//label/textarea").send_keys("test video")
    file = selenium.driver.find_element(By.ID, 'videofile')
    selenium.driver.execute_script("arguments[0].removeAttribute('style')", file)
    selenium.screenshot('file')
    file.send_keys(join(DIR, 'videos', 'test.mp4'))
    publish = "//button[text()='Publish!']"
    selenium.wait_driver.until(EC.element_to_be_clickable((By.XPATH, publish)))
    selenium.find_by_xpath(publish).click()
    selenium.find_by_xpath("//*[text()='test video']")
    selenium.find_by_xpath("//span[text()='New post']")
    assert not selenium.exists_by(By.XPATH, "//span[contains(.,'Error processing')]")
    selenium.screenshot('publish-video')


def test_teardown(driver):
    driver.quit()
