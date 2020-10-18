/*!
 * Level Up
 * Copyright (c) 2020 Alexandra Frock, Cutie Café, contributors
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 * 
 */

export default async function activate(){
    await (async () => {
        let items;
        let password;
        let apps;

        while( true ){
            password = prompt("Enter the admin password, or press enter/click cancel to exit", "");

            if( password == "" ) return;
        
            try {
                items = await (await fetch("/api/admin?key=" + password)).json();
                apps = await (await fetch("/api/suggestions")).json()
            } catch(e){
                continue;
            }

            break;
        }

        while( true ){
            let appid = prompt("Enter an AppID to manipulate, or anything else to exit.\n\nThere are " + Object.keys(items.UnapprovedApps).length + " unapproved apps:\n\n" + Object.keys(items.UnapprovedApps).join("\n") + "", "");

            if( isNaN(parseInt(appid)) ) return;
            
            while( true ) {
                let prmpt = "AppID " + appid + "\n";

                let appType = 0;
        
                if( Object.keys(apps).indexOf(appid) > -1 ){
                    prmpt += "Is an approved, suggested app (" + apps[appid].Name + ")\n"
                    appType = 1;
                }
                else if( Object.keys(items.UnapprovedApps).indexOf(appid) > -1 ){
                    prmpt += "Is an unapproved app\n"
                    appType = 0;
                }
                else {
                    prmpt += "Is a new app\n";
                    appType = -1;
                }

                prmpt += "\n";

                if( items.UnapprovedApps[appid] ) prmpt += (items.UnapprovedApps[appid].Review || "(no review)") + "\n\n";

                if( appType == 0 || appType == 1 ){
                    prmpt += "Enter U to unapprove and delete this app.\n";
                }

                if( appType == -1 || appType == 0 ){
                    prmpt += "Enter A to add and/or approve this app. NOTE: if the server is running in SERVE MODE ONLY, prices may be inaccurate until the fetch bot runs\n";
                }

                prmpt += "Enter O to open this app in a new tab.\n";

                prmpt += "Enter nothing to exit.";
        
                let resp = prompt(prmpt, "");
                if( resp == "" || resp == null ) break;

                switch(resp.toLowerCase()){
                    case "a":
                        try {
                            items = await (await fetch("/api/admin?key=" + password + "&action=approve&appid=" + appid)).json();
                            apps = await (await fetch("/api/suggestions")).json()
                        } catch(e){
                            alert("Action failed: " + e);
                        }
                        break;
                    case "u":
                        try {
                            items = await (await fetch("/api/admin?key=" + password + "&action=delete&appid=" + appid)).json();
                            apps = await (await fetch("/api/suggestions")).json()
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