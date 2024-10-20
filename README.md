# crypto-playground

## axs-restake-reward
Axie Infinity Shard(AXS) auto compounding program.

### Motivation
Checking the AXS staking dashboard manually to see if 24 hours had passed and then requesting to compound can be tedious.

### Running axs-restake-reward
To automatically compund your AXS rewards every day:

1. Create a config folder:
- In the root directory of your project, create a new folder named `config`.

2. Create a `axs_staking_info.yml` file inside the config folder:
- Within the `config` folder, create a file named `axs_staking_info.yml`.

3. Add the following content to the `axs_staking_info.yml` file:
- `chainID: 2020`
- `gasLimit: 371098`
- `pk: {your private key starting without 0x}`
  - Your wallet address will be derived from the private key.

4. This program automatically tracks your previous compound time and compounds your rewards every 24 hours.
- AXS allows you to coumpound every 24hours.
- If 24 hours have passed since the last compound, perform an immediate compound.
- If not, wait until 24 hours have passed since the last compound.

Thereafter, the program will automatically compound every 24 hours + 1 minute.

5. (Optional) You can send Telegram notifications
- You can receive Telegram notifications by setting the flag -telegram to true.

```
./crypto-playground -telegram=true
```
- Also, you need to add telegram configurations, in `/config/axs_staking_info.yml`
    - `telegram: `
      - `token`
      - `chatId`
      - `userName`
      - `webHookUrl`

- You can interact with the bot using certain commands. <br/>
![image-Photoroom (2)](https://github.com/user-attachments/assets/17d80ca6-aca4-4381-93f7-a51638aeb3ec)

### Unlock Schedules for AXS in Detail
![image](https://github.com/user-attachments/assets/5a629a26-6b7a-4906-a80f-10c5c3652ca9)
<img width="612" alt="image" src="https://github.com/user-attachments/assets/53486155-6719-467e-b246-b4524517fdae">
