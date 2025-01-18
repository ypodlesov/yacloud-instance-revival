# Yandex Cloud Instance Revival

This is a service for revival of interruptible compute instances in [Yandex Cloud Service](https://yandex.cloud/en)

### How to use

1. Get IAM token. [How to](https://yandex.cloud/en/docs/iam/operations/iam-token/create).
2. Paste it in file `systemd/env`

```
IAM_TOKEN=<tour_token>
CONFIG_PATH=/etc/yacloud_revival/config.yaml
```

3. Get desired `instance_id` list and place it in `configs/config.yaml`.

3. Run `sudo chmod +x setup.sh && sudo ./setup.sh`

4. To check everything is ok You can run `sudo systemctl status yacloud_revival.service`


### Configuration

In file `configs/config.yaml` You can configure how often the Instance Revival should ping your desired hosts. You can set general value `check_health_period` and override it for each instance.

### Problems

If something went wrong you can view logs in file `/etc/yacloud_revival/general.log`. This file is erased once a day by `log_eraser.service` on `log_eraser.timer`.