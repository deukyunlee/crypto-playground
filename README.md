# crypto-playground

## axs-restake-reward
Automated Axie Infinity Shard(AXS) restaking program.

### Motivation
I found it tedious to check the AXS staking dashboard every time to see if 24 hours had passed and then manually restake. 
So, I thought of automating this process.

### Running axs-restake-reward
To automatically restake your AXS rewards every day:

1. Create a config folder:
- In the root directory of your project, create a new folder named `config`.

2. Create a `axs_staking_info.yml` file inside the config folder:
- Within the `config` folder, create a file named `axs_staking_info.yml`.

3. Add the following content to the `axs_staking_info.yml` file:
- `chainID: 2020`
- `gasLimit: 371098`
- `accountAddress: {your account address starting with 0x}`
- `pk: {your private key starting with 0x}`

4. When running the program, you need to specify the time and minutes for the initial restaking. 
- Since AXS prevents restaking for 24 hours, it is necessary to know the previous restaking time to ensure the automated restaking occurs correctly.
- To handle this, provide the time in RFC3339 format as a flag value when you first run the program. The program will use this time to schedule automatic restaking, which will occur every 24 hours and 1 minute from the specified time.

Example: Suppose you have restaked your reward at 2024-07-20T22:58:16+09:00.

```
./crypto-playground -time=2024-07-20T22:58:16+09:00
```

In this example, the program will execute the initial restaking based on the previous restaking time:

- If 24 hours have passed since the last restaking, perform an immediate restaking.
- If not, wait until 24 hours have passed since the last restaking.

Thereafter, the program will automatically restake every 24 hours and 1 minute.

5. (Optional) You can send Telegram notifications
- If you set flag `-telegram` to `true`, you can receive telegram notifications.

```
./crypto-playground -time=2024-07-20T22:58:16+09:00 -telegram=true
```
- Also, you need to add telegram configurations, in `/config/axs_staking_info.yml`
    - `telegramToken`
    - `telegramChatId`
    - `telegramUserName`

![image](https://github.com/user-attachments/assets/0ff800d9-6843-425c-b799-6d5d6160bd70)
