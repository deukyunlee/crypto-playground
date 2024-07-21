# crypto-playground

## axs-restake-reward
### config
in order to restake your AXS rewards automatically every day,

1. Create a config folder:
- In the root directory of your project, create a new folder named config.

2. Create a axs_staking_info.yml file inside the config folder:
- Within the config folder, create a file named axs_staking_info.yml.

3. Add the following content to the axs_staking_info.yml file:
- chainID: 2020
- gasLimit: 371098
- accountAddress: {your accountAddress starting with 0x}
- pk: {your private key starting with 0x}

<img width="187" alt="image" src="https://github.com/user-attachments/assets/c8c7eec3-dea0-40c4-8dac-50272358922d">

4. When running the program, you need to specify the time and minutes for the initial restaking. 
- Since AXS prevents restaking for 24 hours, it is necessary to know the previous restaking time to ensure the automated restaking occurs correctly.
- To handle this, provide the time as a flag value when you first run the program. The program will use this time to schedule automatic restaking, occurring every 24 hours plus 1 minute from the specified time.

e.g. Suppose you have restaked your reward 02:00 AM yesterday.

```
./crypto-playground -hour=2 -minute=0
```

In this example, the program will schedule the first restaking for 02:01 AM and then automatically restake every 24 hours plus 1 minute thereafter.
