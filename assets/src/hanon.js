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

import $ from 'jquery';
import DOMPurify from 'dompurify';
import activate from './chico';

const days = [ "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday" ];
const months = [ "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December" ]

let apps;
let specialEvent;
let currency = "$";
let country = "us";

let args = {};

let selectedGenre = "";
let underPrice = -1;
let sortType = "name";
let os = ""
let demo = ""
let discounted = ""
let page = 1

let appsPerPage = 10;

let appSearchHTML;

let genres = {};
let genreSections = [
  [
    "Action",
    "Adventure",
    "Casual",
    "Indie",
    "Racing",
    "RPG",
    "Simulation",
    "Sports",
    "Strategy",
    "Metroidvania",
    "Sandbox",
    "Puzzle",
    "Visual Novel",
    "FPS",
    "Arcade",
    "Platformer",
    "Shooter",
    "Souls-like",
    "Survival",
    "City Builder",
    "Interactive Fiction",
    "Twin Stick Shooter",
    "3D Platformer",
    "Base Building",
    "Side Scroller",
    "Third-Person Shooter",
    "Puzzle Platformer",
    "Beat 'em Up",
    "Roguelike",
    "Time Management",
    "Trading Card Game",
    "Shoot 'Em Up",
    "Turn-Based"
  ],
  [
    "Singleplayer",
    "Multiplayer",
    "Online Co-Op",
    "Local Multiplayer",
    "Massively Multiplayer",
    "Local Co-Op",
    "Co-op",
    "Split Screen"
  ]
]

let appReviewHelpful = ( window.localStorage.appReviewHelpful ? window.localStorage.appReviewHelpful.split(",") : [] );

// -- utility functions
function clone(elem){
  let nelem = elem.cloneNode(true);
  elem.parentNode.replaceChild(nelem, elem);
  return nelem;
}

function formatPrice(price, forceUS){
  if( country == "fr" && ! forceUS ) price = price.toString().replace(".", ",") + "&euro;"
  else if( country == "uk" && ! forceUS ) price = "&pound;" + price.toString();
  else price = "$" + price.toString() + " USD";

  return price;
}

function replaceHashParam(param, wth){
  args[param] = wth;
  if( param != "page" ) args.page = 1;
  window.location.hash = Object.keys(args).filter(i => i.length > 0).map(i => i + "=" + args[i]).join(";");
}

function buyAppButton(app, w100 = false){
  return `${app.price == null ? `unavailable` : `<a class="btn btn-sm btn-primary${w100 ? " w-100" : ""}" target="_blank" href="${app.price.url}">${app.price.provider} (${app.price.discount > 0 ? "-" + app.price.discount + "% " : ""}${app.price.price > 0 ? formatPrice(app.price.price/100, app.price.provider == "Humble") : "Free"})</a><br>`}`;
}

function ratingSpan(app){
  let c = "None";

  if( app.Score*100 >= 80 ) c = "Positive"
  else if( app.Score*100 >= 40 ) c = "Mixed";
  else c = "Negative";

  return `<span class="${c.toLowerCase()}" title="Score: ${Math.floor(app.Score*100)}%">${c}</span>`;
}

function addPaginator(page, maxPages){
  let html = "";

  if( page > 1 ) html += `<p style="float: left;"><a class="replace" data-arg="page" data-replace="${parseInt(page)-1}" href="#">Previous</a></p>`;
  if( page < Math.ceil(maxPages) ) html += `<p style="float: right;"><a class="replace" data-arg="page" data-replace="${parseInt(page)+1}" href="#">Next</a></p>`;

  html += `<p style="text-align: center;">Page <b>${page}</b> of ${Math.ceil(maxPages)}</p>`;

  return html
}

// -- interactable functions
async function wasAppHelpful(appid, wasHelpful){
  if( appReviewHelpful.indexOf(appid) > 0 ) return false;
  appReviewHelpful.push(appid);

  let resp = await (fetch("/api/helpful", {
    method: "POST",
    body: JSON.stringify({
        AppID: appid,
        WasHelpful: wasHelpful,
    }),
    headers: {
        "Content-Type": "application/json"
    }
  }))

  window.localStorage.appReviewHelpful = appReviewHelpful;

  if( resp.status != 200 ){
    alert(await resp.text())
    return false;
  }

  return true;
}

function scanPrices(){
  for( let app in apps ){
    app = apps[app];

    app.price = app.Prices[Object.keys(app.Prices).sort((j, k) => {
      if( app.Prices[j] && app.Prices[j][country] ) app.Prices[j][country].provider = j;
      if( app.Prices[k] && app.Prices[k][country] ) app.Prices[k][country].provider = k;

      if( ! app.Prices[j] || ! app.Prices[j][country] ) return 1;
      if( ! app.Prices[k] || ! app.Prices[k][country] ) return -1;

      if( app.Prices[j][country].price < app.Prices[k][country].price ) return -1;
      else if( app.Prices[j][country].price > app.Prices[k][country].price ) return 1;
      else return 0;
    })[0]];

    if( app.price != null ) app.price = app.price[country];
  }
}

function getAppFromSearchBox(val){
  if( val == "" ) return 0;

  if( isNaN(parseInt(val)) ){
    let store = /^(https?:\/\/)?(store\.steampowered\.com|steamcommunity\.com)\/app\/(\d{1,8})/.exec(val);
    if( store && store.length > 3 ){
      return parseInt(store[3]);
    }
  } else {
    return parseInt(val);
  }
}

function initSearch(){
  $("#submit").on("click", () => {
    let app = getAppFromSearchBox($("#appsearch").val());

    if( app === 0 ) return $("#error").text("Please enter the Store page of the title you want to suggest.");
    if( $("#review").val().trim() == "" || $("#review").val().trim().length < 10 || $("#review").val().trim().length > 300 ) return $("#error").text("You must include why you want to recommend this title, and your response must be 10 to 300 characters long.");
    if( grecaptcha.getResponse() == "" ) return $("#error").text("Are you a robot?");

    $("#error").text("");
    $("#submit").addClass("disabled");
    $("#submit").text("Please wait...");

    fetch("/api/suggestions", {
      method: "POST",
      body: JSON.stringify({
          AppID: app,
          Review: $("#review").val(),
          Recaptcha: grecaptcha.getResponse()
      }),
      headers: {
          "Content-Type": "application/json"
      }
    }).then((resp) => {
      if( resp.status != 200 ){
        resp.text().then((x) => {
          $("#error").text(x);
          $("#submit").removeClass("disabled");
          $("#submit").text("Submit");
        })
        return;
      }

      $("#submit-modal").find(".modal-body").html("Thanks for your submission! <a id='submit-another' href='#'>Submit another?</a>")

      $("#submit-another").on("click", () => {
        document.querySelector("#submit-modal").querySelector(".modal-body").innerHTML = appSearchHTML;
        $("#submit").removeClass("disabled");
        $("#submit").text("Submit");
        $("#submit").attr("style", "display: block;");
        initSearch();
        grecaptcha.render(document.querySelector(".g-recaptcha"));
      });

      $("#submit").attr("style", "display: none;");
    }).catch(() => {
      $("#error").text("Could not submit. Try again later.");
      $("#submit").removeClass("disabled");
      $("#submit").text("Submit");
    });
  });

  $("#show-submit-modal").on("click", e => {
    e.preventDefault();
    $("#submit-modal").modal('show');
  })

  $("#show-store-tags-modal").on("click", e => {
    e.preventDefault();
    $("#store-tags-modal").modal('show');
  });
}

function refreshApps(lApps, page = 1, maxPages = 1){
  window.scrollTo(0, 0);

  Array.from(document.querySelectorAll(".currency")).forEach(i => i.innerHTML = currency);

  if( sortType != "added_asc" || selectedGenre != "" || underPrice != -1 || os != "" || demo != "" || page > 1 || discounted != "" ){
    document.querySelector("#app-carousel").setAttribute("style", "display: none");
    document.querySelector("#apps").setAttribute("style", "margin-top: calc(56px + .5em);")
  } else {
    document.querySelector("#app-carousel").setAttribute("style", "display: block");
    document.querySelector("#apps").setAttribute("style", "");

    let chosen = [];

    Array.from(document.querySelectorAll("[data-lu-slide]")).forEach(i => {
      let choice = -1;

      while( chosen.indexOf(choice) > -1 || choice == -1 ) {
        choice = Object.keys(apps)[Math.floor(Math.random()*Object.keys(apps).length)];
        if( Object.keys(apps).length < 5 ) break;
      }

      let app = apps[choice];
      i.querySelector("img").setAttribute("src", app.Screenshot);
      i.querySelector("h5").innerText = app.Name;
      i.querySelector("p").innerHTML = DOMPurify.sanitize(app.Description) + "<br><br>" + buyAppButton(app);

      chosen.push(choice);
    });
  }

  let html = "";

  if( lApps.length < 1 ){
    html = "<h1>Oops!</h1>Your search didn't turn up anything. Broaden your search terms to find something you'll love.";  
  } else {

    html += `<div id="special-event">${putSpecialEvent()}</div>`;

    html += addPaginator(page, maxPages);

    html += `<div class="list-group">`;

    for( let app in lApps ){
      app = lApps[app];

      html += `
        <div class="list-group-item flex-cloumn align-items-start">
          <div class="d-flex w-100">
            <img class="app-picture" data-appid="${app.AppID}" src="https://cdn.cloudflare.steamstatic.com/steam/apps/${app.AppID}/capsule_184x69.jpg">
            <div class="ml-1 flex-fill">
              <h5 class="mb-1 mt-0">${DOMPurify.sanitize(app.Name)}</h5>
              <h6 class="mb-1 text-muted">${DOMPurify.sanitize(app.Developers[0].trim() == app.Publishers[0].trim() ? app.Developers[0] : app.Developers[0] + "; " + app.Publishers[0])}</h6>
            </div>

            <div class="mb-1 d-none d-md-block">
              <p style="text-align: center;">
                ${buyAppButton(app)}
              </p>
            </div>
          </div>

          <p class="mb-0">
            <small class="text-muted">
              ${app.Platforms.Windows ? `<span class="platform windows" title="Windows">&nbsp;</span>`: ""}
              ${app.Platforms.MacOS ? `<span class="platform mac" title="macOS">&nbsp;</span>`: ""}
              ${app.Platforms.Linux ? `<span class="platform linux" title="SteamOS/Linux">&nbsp;</span>`: ""}
              ${ratingSpan(app)}
               &middot; <a href="https://s.team/a/${app.AppID}">view on Steam</a>
              ${app.Demo ? ` &middot; <a href="https://store.steampowered.com/app/${app.AppID}/#game_area_purchase" target="_blank">demo available</a>` : ""}
            </small>
            
            <br>

            ${app.Genres.filter(i => i != "Early Access").slice(0, 10).map(i => "<a href='#' class='badge badge-info tag replace' data-arg='genre' data-replace='" + DOMPurify.sanitize(i) + "'>" + DOMPurify.sanitize(i) + "</a>").join(" ")}
          </p>

          <p class="mb-0">
            ${DOMPurify.sanitize(app.Description)}
          </p>

          ${app.Review ? `
            <blockquote class="review blockquote">
              ${DOMPurify.sanitize(app.Review)}
              <footer class="blockquote-footer"><cite>the recommender</cite></footer>
            </blockquote>
            
            <p class="text-muted mb-0" style="font-size: 85%">
              ${appReviewHelpful.indexOf("" + app.AppID) < 0 ? `
                Was this recommendation helpful? <a class="badge badge-primary helpful" href="#" data-recommend="true" data-appid="${app.AppID}">Yes</a> <a class="badge badge-danger helpful" href="#" data-recommend="false" data-appid="${app.AppID}">No</a><br>
              ` : ""}

              ${app.HelpfulTotal > 0 ? `
                ${app.HelpfulPositive} ${app.HelpfulPositive == 1 ? "person" : "people"} found this helpful
              ` : ""}
            </p>
          ` : ""}

          <div class="mt-2 d-xs-block d-sm-block d-md-none w-100">
            <p style="text-align: center;">
              ${buyAppButton(app, true)}
            </p>
          </div>
        </div>
      `;
    }

    html += "</div>";

    html += addPaginator(page, maxPages);
  }

  document.querySelector("#apps").innerHTML = html;

  document.querySelectorAll(".replace").forEach(i => {
    i = clone(i);

    i.addEventListener("click", e => {
      e.preventDefault();
      replaceHashParam(i.getAttribute("data-arg"), i.getAttribute("data-replace"));
    });
  });

  document.querySelectorAll(".fprice").forEach(i => {
    i.innerHTML = formatPrice(i.getAttribute("data-price"));
  });

  document.querySelectorAll(".helpful").forEach(i => {
    i.addEventListener("click", async e => {
      e.preventDefault();

      let pe = i.parentElement;
      pe.innerText = "working...";
      if( ! await wasAppHelpful(parseInt(i.getAttribute("data-appid")), i.getAttribute("data-recommend") == "true") ) pe.innerText = "error";
      else pe.innerText = "Thanks for your feedback!";
    });
  });
}

function parseHash(){
  let params = {};

  window.location.hash.substring(1).replace(/\%20/g, " ").split(";").forEach(i => {
    params[i.split("=")[0]] = i.split("=")[1];
  });

  args = params;

  if(params.admin == "1"){
    activate();
    return;
  }

  selectedGenre = params.genre || "";

  sortType = params.sort || "added_asc";

  os = params.os || "";

  demo = params.demo || "";

  page = params.page || 1;

  country = params.cc || "us";
  if( country == "fr" ) currency = "&euro;";
  else if( country == "uk" ) currency = "&pound;";
  else currency = "$";

  underPrice = ! isNaN(parseFloat(params.under)) ? parseFloat(params.under) : -1;

  discounted = params.discounted || "";

  scanPrices();

  let apply = Object.keys(apps).filter(i => {
    if( selectedGenre != "" ){
      if( apps[i].Genres.map(i => i.toLowerCase()).indexOf(selectedGenre.toLowerCase()) < 0 ) return false;
    }

    if( underPrice > 0 ){
      if( (apps[i].price.price/100) > underPrice ) return false;
    }

    if( os != "" ){
      if( os == "macos" && ! apps[i].Platforms.MacOS ) return false;
      else if( ( os == "linux" || os == "steamos" ) && ! apps[i].Platforms.Linux ) return false;
      else if( os == "windows" && ! apps[i].Platforms.Windows ) return false;
    }

    if( demo != "" && ! apps[i].Demo ) return false;

    if( discounted != "" && apps[i].price.discount < 1 ) return false;

    return true;
  }).map(i => apps[i]);

  switch(sortType){
    case "old":
      break;
    case "new":
      apply = apply.reverse();
      break;
    case "price_asc":
      apply = apply.sort((a, b) => {
        if( a.price.price < b.price.price ) return -1;
        else if( b.price.price < a.price.price ) return 1;
        else return 0;
      });
      break;
    case "price_desc":
      apply = apply.sort((a, b) => {
        if( a.price.price > b.price.price ) return -1;
        else if( b.price.price < a.price.price ) return 1;
        else return 0;
      });
      break;
    case "wilson":
    default:
      apply = apply.sort((a, b) => {
        if( a.Score < b.Score ) return 1;
        else if( a.Score > b.Score ) return -1;
        else return 0;
      });
  }

  let pages = apply.length/appsPerPage;

  apply = apply.slice(appsPerPage*(page-1), appsPerPage*page);

  refreshApps(apply, page, pages);
}

function applyTheme(){
  if( window.localStorage.dark === "true" ){
    $("body").attr("class", "dark");
  } else {
    $("body").attr("class", "");
  }
}

// -- hash change
window.addEventListener("hashchange", parseHash);

window.addEventListener("load", async () => {
  $("#light").on("click", e => {
    e.preventDefault();

    window.localStorage.dark = false;
    applyTheme();
  });

  $("#dark").on("click", e => {
    e.preventDefault();

    window.localStorage.dark = true;
    applyTheme();
  });

  applyTheme();
});

function pad(n){
  return (n < 10 ? "0" : "" ) + n;
}

function putSpecialEvent(){
  if( specialEvent == null ) return "";

  let niceTimeToEnd = "";

  let endTime = new Date(specialEvent.EndTimestamp*1000);

  let timeToEnd = endTime - new Date();
  if( timeToEnd <= 0 ) return "";

  timeToEnd/=1000;

  if( timeToEnd < 60 ) niceTimeToEnd = "in less than a minute";
  else if( timeToEnd < 3600 ) niceTimeToEnd = `in about ${Math.round(timeToEnd/60)} minutes`;
  else if( timeToEnd < 86400 ) niceTimeToEnd = `in about ${Math.round(timeToEnd/3600)} hours`;
  else niceTimeToEnd = `in about ${Math.round(timeToEnd/86400)} days`;

  let niceTimestamp = days[endTime.getDay()] + " " + months[endTime.getMonth()] + " " + endTime.getDate() + " " + endTime.getFullYear() + ", " + (endTime.getHours()%12) + ":" + pad(endTime.getMinutes()) + " " + (endTime.getHours() > 12 ? "pm" : "am");

  return `The ${specialEvent.Name} ends <time title="${niceTimestamp}">${niceTimeToEnd}</time>! <a href="#discounted=1">${Object.keys(apps).filter(i => apps[i].price.discount > 0).length} recommendations are discounted.</a>`;
}

// -- load
window.addEventListener("load", async () => {
  let newDiv = document.createElement("div");

  newDiv.innerHTML = document.querySelector("#submit-modal > .modal-dialog > .modal-content > .modal-body").innerHTML;
  newDiv.querySelector(".g-recaptcha").innerHTML = "";

  appSearchHTML = newDiv.innerHTML;

  let se = false;

  document.querySelector("#send-feedback").addEventListener("click", e => {
    if( ! se ){
      e.preventDefault();
      let dev = "cafe.eituc@srepoleved".split("").reverse().join("").replace("efac", "cafe");
      e.target.innerText =  "send feedback to " + dev;
      e.target.setAttribute("href", "mailto:" + dev);
      se = true;
    }
  });

  let response = await (await fetch("/api/data")).json();
  apps = response.Apps;
  specialEvent = response.Event;

  for(let app of Object.keys(apps)){
    apps[app].Genres.forEach(i => {
      if( ! genres[i] ) genres[i] = 0;
      genres[i]++;
    });
  }

  let sortedGenres = Object.keys(genres).sort((a,b) => {
    if(genres[a] > genres[b]) return -1;
    else if( genres[a] < genres[b] ) return 1;
    else {
      if( a < b ) return -1;
      else if( a > b ) return 1;
      else return 0;
    }
  });

  let sections = [];
  for( let i=0; i<genreSections.length; i++ ){
    sections[i] = [];

    for( let j of sortedGenres ){
      if( genreSections[i].indexOf(j) > -1 ){
        sections[i].push(j);
      }
    }

    for( let j of sections[i] ){
      sortedGenres.splice(sortedGenres.indexOf(j), 1);
    }
  }

  sections[sections.length] = sortedGenres;

  let genresHtml = "";
  for(let i in sections ){
    for( let j of sections[i] ){
      genresHtml += `
        <a class="badge badge-info replace" href="#" data-arg="genre" data-replace="${j}" data-dismiss="modal">${j} (${genres[j]})</a>
      `;
    }
    genresHtml += "<hr>";
  }
  document.querySelector("#store-tags-modal-body").innerHTML = genresHtml;

  parseHash();
  initSearch();
});