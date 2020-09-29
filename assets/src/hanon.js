import $ from 'jquery';
import DOMPurify from 'dompurify';
import activate from './chico';

let apps;
let currency = "$";
let country = "us";

let args = {};

let selectedGenre = "";
let underPrice = -1;
let sortType = "name";
let os = ""
let demo = ""
let page = 1

let appsPerPage = 10;

let genres = [];

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
  window.location.hash = Object.keys(args).filter(i => i.length > 0).map(i => i + "=" + args[i]).join(";");
}

function buyAppButton(app){
  return `${app.price == null ? `unavailable` : `<a class="btn btn-sm btn-primary" target="_blank" href="${app.price.url}">${app.price.provider} (${app.price.discount > 0 ? "-" + app.price.discount + "% " : ""}${app.price.price > 0 ? formatPrice(app.price.price/100, app.price.provider == "Humble") : "Free"})</a>`}`
}

function addPaginator(page, maxPages){
  let html = "";

  if( page > 1 ) html += `<p style="float: left;"><a class="replace" data-arg="page" data-replace="${parseInt(page)-1}" href="#">Previous</a></p>`;
  if( page < Math.ceil(maxPages) ) html += `<p style="float: right;"><a class="replace" data-arg="page" data-replace="${parseInt(page)+1}" href="#">Next</a></p>`;

  html += `<p style="text-align: center;">Page <b>${page}</b> of ${Math.ceil(maxPages)}</p>`;

  return html
}

// -- interactable functions
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

function initSearch(){
  $("#appsearch").typeahead({
    order: "asc",
    dynamic: true,
    source: {
      app: {
        display: "name",
        template: (_, item) => {
          return item.name + " (" + item.appid + ")"; 
        },
        data: [],
        ajax: {
          type: "GET",
          url: "/api/search?q={{query}}"
        }
      }
    },
    callback: {
      onClick: (node, a, item, event) => {
        event.preventDefault();
        document.querySelector("#appsearch").value = item.name + " (" + item.appid + ")";
      }
    }
  });

  $("#submit").on("click", () => {
    if( $("#appsearch").val() == "" || ( isNaN(parseInt($("#appsearch").val())) && /\((\d*)\)$/.exec($("#appsearch").val()) == null ) ) return $("#error").text("Please enter an app to suggest.");
    if( grecaptcha.getResponse() == "" ) return $("#error").text("Are you a robot?");

    $("#error").text("");
    $("#submit").addClass("disabled");
    $("#submit").text("Please wait...");

    fetch("/api/suggestions", {
      method: "POST",
      body: JSON.stringify({
          AppID: isNaN(parseInt($("#appsearch").val())) ? parseInt(/\((\d*)\)$/.exec($("#appsearch").val())[1]) : (parseInt($("#appsearch").val())),
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

      $(".modal-body").text("Thanks for your submission!")
      $("#submit").attr("style", "display: none;");
    }).catch(() => {
      $("#error").text("Could not submit. Try again later.");
      $("#submit").removeClass("disabled");
      $("#submit").text("Submit");
    });
  });

  $("#show-submit-modal").on("click", e => {
    e.preventDefault();
    $("#submit-modal").modal('show')
  })
}

function refreshApps(apps, page = 1, maxPages = 1){
  Array.from(document.querySelectorAll(".currency")).forEach(i => i.innerHTML = currency);

  if( sortType != "added_asc" || selectedGenre != "" || underPrice != -1 || os != "" || demo != "" || page > 1 ){
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

  if( apps.length < 1 ){
    html = "<h1>Oops!</h1>Your search didn't turn up anything. Broaden your search terms to find something you'll love.";  
  } else {

    html += addPaginator(page, maxPages);

    html += `<div class="list-group">`;

    for( let app in apps ){
      app = apps[app];

      html += `
        <div class="list-group-item flex-cloumn align-items-start">
          <div class="d-flex w-100">
            <img class="app-picture" data-appid="${app.AppID}" src="https://cdn.cloudflare.steamstatic.com/steam/apps/${app.AppID}/capsule_184x69.jpg">
            <div class="ml-1 flex-fill">
              <h5 class="mb-1 mt-0">${DOMPurify.sanitize(app.Name)}</h5>
              <h6 class="mb-1 text-muted">${DOMPurify.sanitize(app.Developers[0].trim() == app.Publishers[0].trim() ? app.Developers[0] : app.Developers[0] + "; " + app.Publishers[0])}</h6>
            </div>

            <div class="mb-1">
              <p style="text-align: center;">
                ${buyAppButton(app)}
              </p>
            </div>
          </div>

          <p class="mb-0">
            ${app.Genres.filter(i => i != "Early Access").map(i => "<a href='#' class='badge badge-info tag replace' data-arg='genre' data-replace='" + DOMPurify.sanitize(i) + "'>" + DOMPurify.sanitize(i) + "</a>").join(" ")}
          </p>

          <p class="mb-0">
            ${DOMPurify.sanitize(app.Description)}
          </p>

          <p class="mb-0">
            <small class="text-muted">
              ${app.Platforms.Windows ? `<span class="platform windows" title="Windows">&nbsp;</span>`: ""}
              ${app.Platforms.MacOS ? `<span class="platform mac" title="macOS">&nbsp;</span>`: ""}
              ${app.Platforms.Linux ? `<span class="platform linux" title="SteamOS/Linux">&nbsp;</span>`: ""}
              ${app.Demo ? `<a href="https://store.steampowered.com/app/427520/Factorio/#game_area_purchase" target="_blank">demo available</a>` : ""}
            </small>
          </p>
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

    return true;
  }).map(i => apps[i]);

  scanPrices();

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

// -- hash change
window.addEventListener("hashchange", parseHash);

// -- load
window.addEventListener("load", async () => {
  apps = await (await fetch("/api/suggestions")).json();
  console.log(apps);

  for(let app of Object.keys(apps)){
    apps[app].Genres.forEach(i => {
      if( genres.indexOf(i) < 0 ){
        genres.push(i);
      }
    });
  }

  let genresHtml = "";
  for(let i of genres ){
    genresHtml += `
      <a class="dropdown-item replace" href="#" data-arg="genre" data-replace="${i}">${i}</a>
    `;
  }
  document.querySelector("#genre-dropdown-options").innerHTML = genresHtml;

  parseHash();
  initSearch();
});