async function activate(){
    await (async () => {
        let items;
        let password;

        while( true ){
            password = prompt("Enter the admin password, or press enter/click cancel to exit", "");

            if( password == "" ) return;
        
            try {
                items = await (await fetch("/api/admin?key=" + password)).json();
            } catch(e){
                continue;
            }

            break;
        }

        while( true ){
            let appid = prompt("Enter an AppID to manipulate, or anything else to exit.\n\nThere are " + items.UnapprovedApps.length + " unapproved apps:\n\n" + items.UnapprovedApps.join("\n") + "", "");

            if( isNaN(parseInt(appid)) ) return;
            
            while( true ) {
                let prmpt = "AppID " + appid + "\n";

                let appType = 0;
        
                if( Object.keys(apps).indexOf(appid) > -1 ){
                    prmpt += "Is an approved, suggested app (" + apps[appid].Name + ")\n"
                    appType = 1;
                }
                else if( items.UnapprovedApps.indexOf(parseInt(appid)) > -1 ){
                    prmpt += "Is an unapproved app\n"
                    appType = 0;
                }
                else {
                    prmpt += "Is a new app\n";
                    appType = -1;
                }

                console.log(items.UnapprovedApps);

                prmpt += "\n";

                if( appType == 0 || appType == 1 ){
                    prmpt += "Enter U to unapprove and delete this app.\n";
                }

                if( appType == -1 || appType == 0 ){
                    prmpt += "Enter A to add and/or approve this app.\n";
                }

                prmpt += "Enter O to open this app in a new tab.\n";

                prmpt += "Enter nothing to exit.";
        
                let resp = prompt(prmpt, "");
                if( resp == "" || resp == null ) break;

                switch(resp.toLowerCase()){
                    case "a":
                        try {
                            items = await (await fetch("/api/admin?key=" + password + "&action=approve&appid=" + appid)).json();
                        } catch(e){
                            alert("Action failed: " + e);
                        }
                        break;
                    case "u":
                        try {
                            items = await (await fetch("/api/admin?key=" + password + "&action=delete&appid=" + appid)).json();
                        } catch(e) {
                            alert("Action failed: " + e);
                        }
                        break;
                    case "o":
                        window.open("https://store.steampowered.com/app/" + appid);
                        break;
                }
            }
        }
    })();

    window.location.hash = "";
    history.go(0);
}