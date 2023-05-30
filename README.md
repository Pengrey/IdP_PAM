# IdP_PAM

## Project Description
This project is a PAM module that allows users to authenticate using external Identity Providers (IdPs). The module uses OAuth 2.0 with the Device Authorization Grant strategy to authenticate users.
![image](https://github.com/Pengrey/IdP_PAM/assets/55480558/2c615691-ff64-4d13-9fa8-8f788386eb54)

## Project Structure
![image](https://user-images.githubusercontent.com/55480558/236838210-94a208ed-cb80-4018-b231-7d62e018d949.png)

## idp_login
`idp_login` is a command-line application to manage Identity Providers (IdPs) and identity attributes for users in a protected repository.

### Usage
```bash
idp_login <action> [options]
```
- `<action>`: The action to perform, such as setting, changing, deleting, or listing IdPs and identity attributes.
- `[options]`: Optional arguments to specify additional information for the action.

### Actions and Options

1. **Set, change, or delete IdPs and their operational parameters (for host administrators):**

```bash
idp_login manage-idp [--operation set|change|delete] [--idp IDP_NAME] [--params PARAMS]
```

- `--operation`: The operation to perform, e.g. set, change, or delete.
- `--idp IDP_NAME`: The name of the IdP to be managed (required for set and change operations).
- `--params PARAMS`: The operational parameters for the IdP (required for set and change operations).

2. **Set, change, or delete identity attributes for a given IdP for the current user:**

```bash
idp_login manage-attributes [--operation set|change|delete] [--idp IDP_NAME] [--attributes ATTRIBUTES]
```

- `--operation`: The operation to perform, e.g. set, change, or delete.
- `--idp IDP_NAME`: The name of the IdP whose attributes need to be managed (required for set and change operations).
- `--attributes ATTRIBUTES`: The identity attributes for the IdP (required for set and change operations).

3. **List all users with registered IdPs (for host administrators):**

```bash
idp_login list-users
```

4. **List the IdPs registered for the current user and the identity parameters for each IdP:**

```bash
idp_login list-idps
```

## Examples

- To set an IdP with its operational parameters:

```bash
sudo idp_login manage-idp --operation set --idp google --params '{"request_url":"https://accounts.google.com/o/oauth2/device/code","request_arguments":{"client_id":"","scope":""},"user_url":"https://accounts.google.com/o/oauth2/device/usercode","poll_url":"https://accounts.google.com/o/oauth2/token","poll_arguments":{"client_id":"","client_secret":"","device_code":"","grant_type":""}}'
```

```bash
sudo idp_login manage-idp --operation set --idp github --params '{"request_url":"https://github.com/login/device/code","request_arguments":{"client_id":"","scope":""},"user_url":"https://github.com/login/device","poll_url":"https://github.com/login/oauth/access_token","poll_arguments":{"client_id":"","device_code":"","grant_type":""}}'
```

- To change an IdP's operational parameters:

```bash
sudo idp_login manage-idp --operation change --idp google --params '{"request_url":"https://accounts.google.com/o/oauth2/device/code","request_arguments":{"client_id":"","scope":""},"user_url":"https://accounts.google.com/o/oauth2/device/usercode","poll_url":"https://accounts.google.com/o/oauth2/token","poll_arguments":{"client_id":"","client_secret":"","device_code":"","grant_type":""}}'
```

```bash
sudo idp_login manage-idp --operation change --idp github --params '{"request_url":"https://github.com/login/device/code","request_arguments":{"client_id":"","scope":""},"user_url":"","poll_url":"https://github.com/login/oauth/access_token","poll_arguments":{"client_id":"","device_code":"","grant_type":""}}'
```

- To delete an IdP:

```bash
sudo idp_login manage-idp --operation delete --idp google
```

```bash
sudo idp_login manage-idp --operation delete --idp github
```

- To Show the available IdPs:

```bash
idp_login manage-attributes --operation set
```

- To Show the available attributes for an IdP:

```bash
idp_login manage-attributes --operation set --idp google
```

```bash
idp_login manage-attributes --operation set --idp github
```

- To set identity attributes for an IdP:

```bash
idp_login manage-attributes --operation set --idp google --attributes '{"request_url":"https://accounts.google.com/o/oauth2/device/code","request_arguments":{"client_id":"[REDACTED]","scope":"https://www.googleapis.com/auth/userinfo.email"},"user_url":"https://accounts.google.com/o/oauth2/device/usercode","poll_url":"https://accounts.google.com/o/oauth2/token","poll_arguments":{"client_id":"[REDACTED]","client_secret":"[REDACTED]","device_code":"","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}}'
```

```bash
idp_login manage-attributes --operation set --idp github --attributes '{"request_url":"https://github.com/login/device/code","request_arguments":{"client_id":"[REDACTED]","scope":"user:email"},"user_url":"https://github.com/login/device","poll_url":"https://github.com/login/oauth/access_token","poll_arguments":{"client_id":"[REDACTED]","device_code":"","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}}'
```

- To change identity attributes for an IdP:

```bash
idp_login manage-attributes --operation change --idp google --attributes '{"request_url":"https://accounts.google.com/o/oauth2/device/code","request_arguments":{"client_id":"[REDACTED]","scope":"https://www.googleapis.com/auth/userinfo.email"},"user_url":"https://accounts.google.com/o/oauth2/device/usercode","poll_url":"https://accounts.google.com/o/oauth2/token","poll_arguments":{"client_id":"[REDACTED]","client_secret":"[REDACTED]","device_code":"","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}}'
```

```bash
idp_login manage-attributes --operation change --idp github --attributes '{"request_url":"https://github.com/login/device/code","request_arguments":{"client_id":"[REDACTED]","scope":"user:email"},"user_url":"https://github.com/login/device","poll_url":"https://github.com/login/oauth/access_token","poll_arguments":{"client_id":"[REDACTED]","device_code":"","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}}'
```

- To delete identity attributes for an IdP:

```bash
idp_login manage-attributes --operation delete --idp google
```

```bash
idp_login manage-attributes --operation delete --idp github
```


- To list all users with registered IdPs:

```bash
sudo idp_login list-users
```

- To list the IdPs registered for the current user:

```bash
idp_login list-idps
```

## Project Setup
### Pam Python Installation
To install the Pam Python module, run the following command:
```bash
sudo apt-get install libpam0g libpam-runtime libpam0g-dev python2 python2-dev
git clone https://github.com/Ralnoc/pam-python.git
cd pam-python
vim ./src/setup.py # Change the python version to 2.7
vim ./src/test.py  # Change the python version to 2.7
sudo make && make install
ls /lib/security   # Check if pam_python.so is present
```

### Pam Python Module
To install the project, run the following command:
```bash
git clone https://github.com/Pengrey/IAA.git
cd Project_1
sudo python pip install -r ./src/pam/requirements.txt
```

Edit /etc/pam.d/common-auth and add this line to the top of the file:
```bash
auth sufficient pam_python.so /path/to/Project_1/src/login.py
```

Edit /etc/pam.d/common-session and add this line to the top of the file:
```bash
session	sufficient pam_python.so /path/to/Project_1/src/login.py
```

### IdP Login
To install the IdP Login command-line application, run the following command:
```bash
sudo groupadd idpadmins             // Create a group for IdP administrators
sudo usermod -a -G idpadmins root   // Add root to the group
cd Project_1/src/idp_login
go build -o idp_login main.go       // Build the application
sudo chown root idp_login           // Set the owner to root
sudo chmod 4750 idp_login           // Set the setuid bit
sudo mv idp_login /usr/bin          // Move the application to /usr/bin
```
