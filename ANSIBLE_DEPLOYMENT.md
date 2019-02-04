# Ansible Deployment

## Customize it:
* update hosts:
   * vars explained
      * ```
        hostname=<hostname of the raspberry Pi>
        ssh_user=pi    (this must be pi unless you hacked raspbian appart)
        gopath=/home/pi/go/    (this should not change, used for the go build)
        repo=github.com/RebelIT/gome    (if you forked change this to your repo location)
        application=gome    (used for metrics telegraf application tag)
        hardware=pi    (used for metrics telegraf hardware tag)
        telegraf_db=10.10.10.10    (telegraf output)
        redis_host=127.0.0.1    (IP if you redis server. configure redis proper if using external redis)
        redis_memory=50M    (this should be fine for the size of this DB and usage. pi does not have much to spare)
        timezone=America/Chicago    (your local time zone. important if scheduling devices and metrics)
        push_blank_device_list=false    (on app deploy push out a brand new devices list, use on initial deploy only)
        ```
   
   * host /main
      * ```
        [gome:children]
          main    (if you deploy multiple instances add more children )
        
        [main]
          10.0.0.106    (IP or fqdn of the pi to deploy this to)
        ```
        
   * secrets
      * you have to create secrets.yml in `ansible/roles/application/vars/secrets.yml`. Not giving you mine. 
      * ```
        slack_secret: 'xxxxxx/xxxxx/xxxxxx' #slack token    (NOT the full https url, only the workspace/channel/token)
        aws_id: 'xxxxxxxxxxxx'
        aws_secret: 'xxxxxxxxxxxxxxxxxx'
        aws_region: 'us-east-2'
        aws_token: ''
        aws_queue_url: 'https://sqs.us-east-2.amazonaws.com/xxxxxxxxxxx/xxxxxxx'
        rpiot_user: 'user'
        rpiot_token: 'aSyperSecretToken'

        ```
* roles and deploy steps
   * ```
     roles:
       - common    (installs common packages and sets system settings (ntp, hostname, etc...))
       - reboot    (reboots to set some common settings)
       - redis     (installs local redis listening on localhost only)
       - application    (installs this gome application, see below for more details)
     ```
     
     * application role explained:
        * sets up your gopath on the pi
        * does a go get to the repo in the hosts file
        * does a go build for the main package to out the biniary application
        * deploys the two .json config files (devices / secrets) used in the application
        * installs the service in systemd to call the application
        * starts the service which is listening over http with no auth (yet)
        * I've tried a few ways to deploy it and i do like this one the best
        * git tag/branch deployment coming soon (i want this for feature testing)
        
## Run it:

* --ask-sudo-pass may be required if running reboot role due to your local setup
   ```
   ansible-playbook ansible_deploy.yml --ask-vault-pass -i ansible_hosts --ask-sudo-pass
   ```
   
* manually update the devices.json or start adding devices by using the /api/devices endpoint which will 
update redis and the /etc/gome/devices.json. 