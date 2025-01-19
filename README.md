# Yandex Cloud Instance Revival

This is a service for revival of interruptible compute instances in [Yandex Cloud Service](https://yandex.cloud/en)

### How to use

1. Get Oauth token. [How to](https://yandex.cloud/en/docs/iam/concepts/authorization/oauth-token).
2. Paste it in file `systemd/env`

```
YANDEX_OAUTH_TOKEN=<tour_token>
CONFIG_PATH=/etc/yacloud_revival/config.yaml
```

3. Get desired `instance_id` list and place it in `configs/config.yaml`.

4. Run `sudo chmod +x setup.sh && ./setup.sh`

5. To check everything is ok You can run `sudo systemctl status yacloud_revival.service`


### Configuration

In file `configs/config.yaml` You can configure how often the Instance Revival should ping your desired hosts. You can set general value `check_health_period` and override it for each instance.

### Problems

If something went wrong you can view logs in file `/etc/yacloud_revival/general.log`. This file is erased once a day by `log_eraser.service` on `log_eraser.timer`.