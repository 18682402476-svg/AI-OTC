const { ethers } = require("hardhat");

async function main() {
  // 获取部署账户
  const [deployer] = await ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Account balance:", (await deployer.getBalance()).toString());

  // 部署USDC合约
  // console.log("\nDeploying USDC contract...");
  // const USDC = await ethers.getContractFactory("USDC");
  // const usdc = await USDC.deploy(
  //   "USDC",
  //   "USDC",
  //   ethers.utils.parseUnits("1000000", 18), // 初始供应量1,000,000 USDC
  //   deployer.address // 资金池地址设为部署者地址
  // );
  // await usdc.deployed();
  // console.log("USDC contract deployed to:", usdc.address);

  // 部署OTCTrading合约
  console.log("\nDeploying OTCTrading contract...");
  const OTCTrading = await ethers.getContractFactory("OTCTrading");
  const otcTrading = await OTCTrading.deploy(
    deployer.address, // 初始拥有者
    deployer.address, // 手续费接收者
    100 // 手续费比例：1%（100 basis points）
  );
  await otcTrading.deployed();
  console.log("OTCTrading contract deployed to:", otcTrading.address);

  // 部署MON合约
  // console.log("\nDeploying MON contract...");
  // const MON = await ethers.getContractFactory("MON");
  // const mon = await MON.deploy(
  //   "MON",
  //   "MON",
  //   ethers.utils.parseUnits("1000000", 18), // 初始供应量1,000,000 MON
  //   deployer.address // 资金池地址设为部署者地址
  // );
  // await mon.deployed();
  // console.log("MON contract deployed to:", mon.address);

  console.log("\nDeployment completed successfully!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
