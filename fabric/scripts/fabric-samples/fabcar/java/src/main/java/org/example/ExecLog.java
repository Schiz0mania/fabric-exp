package org.example;


import com.google.gson.Gson;
import org.hyperledger.fabric.gateway.*;

import java.nio.file.Path;
import java.nio.file.Paths;

// 用于更新账本，准确的说是上传打包的行为信息
public class ExecLog {
    static class Rinfo{//tag R
        public String pk;
        public String crlinfo;
        public String icode;
        public Rinfo(String a, String b, String c){
            pk=a;
            crlinfo=b;
            icode=c;
        }

    }
    static class Upinfo{//tag U
        public String encrptParam;
        public String CIndex;
        public String Cloc;
        public Upinfo(String a, String b, String c){
            encrptParam=a;
            CIndex=b;
            Cloc=c;
        }

    }
    static class SNFinfo{//tag SF
        public String DoSid;
        public String log;

        public SNFinfo(String a, String b ){
            DoSid=a;
            log=b;
        }

    }
    public static String packRinfo(String pk,String crlinfo,String icode){
        Rinfo upinfo = new Rinfo(pk,crlinfo,icode);
        String actstr =  new Gson().toJson(upinfo);
        return actstr;
    }
    public static String packUpinfo(String encrptParam,String CIndex,String Cloc){
        Upinfo upinfo = new Upinfo(encrptParam,CIndex,Cloc);
        String actstr =  new Gson().toJson(upinfo);
        return actstr;
    }
    public static String packSNFinfo(String DoSid,String log){
        SNFinfo upinfo = new SNFinfo(DoSid,log);
        String actstr =  new Gson().toJson(upinfo);
        return actstr;
    }

    public static void UploadInfo(String key,String sid,String usercrd,String actTag,String actinfo)throws Exception{
        // Load a file system based wallet for managing identities.
        Path walletPath = Paths.get("wallet");
        Wallet wallet = Wallets.newFileSystemWallet(walletPath);
        // load a CCP
        Path networkConfigPath = Paths.get("..", "..", "test-network", "organizations", "peerOrganizations", "org1.example.com", "connection-org1.yaml");

        Gateway.Builder builder = Gateway.createBuilder();
        builder.identity(wallet, "appUser").networkConfig(networkConfigPath).discovery(true);

        // create a gateway connection
        try (Gateway gateway = builder.connect()) {

            // get the network and contract
            Network network = gateway.getNetwork("mychannel");
            Contract contract = network.getContract("MyContract");
            byte[] result;
            contract.submitTransaction("WriteLedger", key,sid,usercrd,actTag,actinfo);

        }

    }


}
