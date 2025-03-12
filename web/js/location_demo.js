/**
 * Location Demo JavaScript
 * Main application logic for the Mist location visualization
 */

// Configuration constants
const API_ENDPOINT = "https://my.locapid.endpoint";
const DEFAULT_MAP_IDX = 0;
const ENTITY_UPDATE_INTERVAL = 2000; // milliseconds
const ZONE_UPDATE_INTERVAL = 2000; // milliseconds
const BLINK_DURATION = 1500; // milliseconds

// Application state
const entityCache = {};
const zoneCache = {};
let map = "";
let mapImage = null;
let mapId = "";
let maps = [];
let mapW = 0;
let mapH = 0;
let zoneStat = null;
let itvlUpdateEntity = null;
let itvlUpdateZone = null;
let isMenuOpen = false;
let markerOpenViaClick = false; // Needed for marker popup behavior

/**
 * Creates a pulsating marker for the map
 * @param {number} radius - Size of the marker in pixels
 * @param {string} color - CSS color for the marker
 * @returns {L.divIcon} Leaflet div icon
 */
const generatePulsatingMarker = function(radius, color) {
    const cssStyle = `
        width: ${radius}px;
        height: ${radius}px;
        background: ${color};
        color: ${color};
        box-shadow: 0 0 0 ${color};
        opacity: 0.6;
    `;

    return L.divIcon({
        html: `<span style="${cssStyle}" class="pulse"/>`,
        className: ''
    });
};

/**
 * Adds blinking effect to an element
 * @param {string} id - DOM element ID
 */
function blink(id) {
    const element = document.getElementById(id);
    if (!element) return;
    
    if (element.classList.contains("blink-highlight")) {
        element.classList.remove("blink-highlight");
        setTimeout(() => {
            element.classList.add("blink-highlight");
        }, BLINK_DURATION);
    } else {
        element.classList.add("blink-highlight");
        setTimeout(() => {
            element.classList.remove("blink-highlight");
        }, BLINK_DURATION);
    }
}

/**
 * Loads map data from the API and initializes the map
 */
function loadMap() {
    const mapApiUrl = `${API_ENDPOINT}/map`;
    
    $.getJSON(mapApiUrl)
        .done((data) => { 
            // Process map data
            for (let i = 0; i < data.length; i++) {
                maps.push({
                    id: data[i].id,
                    height: data[i].height,
                    width: data[i].width,
                    name: data[i].name
                });
            }

            // Set default map if none selected
            if (mapId === "") {
                mapId = data[DEFAULT_MAP_IDX].id;
            }

            // Initialize Leaflet map
            map = L.map("map", {
                crs: L.CRS.Simple,
                maxZoom: 1,
                minZoom: -4,
                zoomControl: false,
                zoomSnap: 0.25,
            });

            // Add zoom control
            L.control.zoom({position: 'bottomright'}).addTo(map);
            
            // Set attribution with timestamp
            const timeNow = new Date().toLocaleString('ja-JP', {timeZone: 'Asia/Tokyo'});
            map.attributionControl.setPrefix(`Last Data Update: <span id="last-update-time">${timeNow}</span> | &copy; Juniper Networks Inc.`);

            // Add Mist logo
            const logoMist = L.control({position: 'bottomleft'});
            logoMist.onAdd = function(map) {
                const div = L.DomUtil.create('div', 'logo-bottomleft');
                div.innerHTML = "<img src='img/juniper_mist.svg' style='width: 180px; opacity: 0.6; margin-left: 7px; margin-bottom: 7px' />";
                return div;
            };
            logoMist.addTo(map);

            // Add Interop logo
	    /*
            const logoInterop = L.control({position: 'topright'});
            logoInterop.onAdd = function(map) {
                const div = L.DomUtil.create('div', 'logo-topright');
                div.innerHTML = "<img src='img/interop24.png' style='width: 240px; opacity: 0.6; margin-right: 7px; margin-top: 7px' />";
                return div;
            };
            logoInterop.addTo(map);
            */

            updateMap();
        })
        .fail((jqXHR, textStatus, errorThrown) => {
            console.error("Failed to load map data:", textStatus, errorThrown);
        });
}

/**
 * Updates the map display with the selected map
 */
function updateMap() {
    if (!mapId) {
        console.error("updateMap was called without mapId set");
        return;
    }

    // Find the selected map data
    let mapdata = null;
    for (let i = 0; i < maps.length; i++) {
        if (maps[i].id === mapId) {
            mapdata = maps[i];
            break;
        }
    }

    if (!mapdata) {
        console.error(`Could not find map with id = ${mapId}`);
        return;
    }

    // Clear existing update intervals
    if (itvlUpdateEntity) {
        clearInterval(itvlUpdateEntity);
        itvlUpdateEntity = null;
    }

    if (itvlUpdateZone) {
        if (zoneStat) {
            zoneStat.remove(map);
            zoneStat = null;
        }
        clearInterval(itvlUpdateZone);
        itvlUpdateZone = null;
    }

    // Set map dimensions
    mapW = mapdata.width;
    mapH = mapdata.height;
    const bounds = [
        [0, 0],           // Southwest corner
        [mapH, mapW]      // Northeast corner
    ];
    
    map.setMaxBounds(bounds);
    map.fitBounds(bounds);

    // Update map image
    const mapImgUri = `img/map/${mapId}.png`;
    
    // Clean up existing map and entities
    if (mapImage && map.hasLayer(mapImage)) {
        // Remove all entity markers
        Object.keys(entityCache).forEach(key => {
            const marker = entityCache[key].marker;
            if (marker) {
                map.removeLayer(marker);
            }
        });
        
        // Clear caches
        Object.keys(entityCache).forEach(key => delete entityCache[key]);
        Object.keys(zoneCache).forEach(key => delete zoneCache[key]);
        
        // Remove map layer
        map.removeLayer(mapImage);
    }

    // Add new map image
    mapImage = L.imageOverlay(mapImgUri, bounds).addTo(map);
    
    // Load zones and entities
    loadZone();
    updateEntity();
    
    // Set update intervals
    itvlUpdateEntity = setInterval(updateEntity, ENTITY_UPDATE_INTERVAL);
}


/**
 * Loads zone data and creates the zone statistics panel
 */
function loadZone() {
    const zoneApiUrl = `${API_ENDPOINT}/zone?map_id=${mapId}`;
    
    $.getJSON(zoneApiUrl)
        .done((data) => {
            zoneStat = L.control({position: 'topright'});
            zoneStat.onAdd = function(map) {
                const div = L.DomUtil.create('div', 'zone-stats');
                let html = `
                  <div class="zone-stat-container">
                    <div class="zone-stat-content">
                      <ul class="zone-stat-content-list">
                        <li>
                          <div class="zone-stat-name" style='justify-content: center; font-weight: bold;'>Zone Name</div>
                          <div class="zone-stat-count" style='font-weight: bold;'>Users</div>
                        </li>
                `;

                // Add each zone to the panel
                for (let i = 0; i < data.length; i++) {
                    if (data[i].map_id !== mapId) {
                        console.warn(`Zone ${data[i].name} has different map_id ${data[i].map_id} than current map_id ${mapId}`);
                        continue;
                    }

                    // Store zone count in cache
                    zoneCache[data[i].id] = data[i].count;
                    
                    // Truncate long zone names
                    let name = data[i].name;
                    if (name.length > 30) {
                        name = name.substring(0, 30) + "...";
                    }

                    html += `
                        <li>
                          <div class="zone-stat-name">${name}</div>
                          <div class="zone-stat-count" id="${data[i].id}">${data[i].count}</div>
                        </li>
                    `;
                }

                html += `
                      </ul>
                    </div>
                  </div>
                `;
                div.innerHTML = html;

                return div;
            };

            zoneStat.addTo(map);
            itvlUpdateZone = setInterval(updateZone, ZONE_UPDATE_INTERVAL);
        })
        .fail((jqXHR, textStatus, errorThrown) => {
            console.error("Failed to load zone data:", textStatus, errorThrown);
        });
}

/**
 * Updates entity data and markers on the map
 */
function updateEntity() {
    const entityApiUrl = `${API_ENDPOINT}/entity`;
    
    $.getJSON(entityApiUrl)
        .done((data) => {
            const newEntityCache = {};
            
            // Process each entity
            for (let i = 0; i < data.length; i++) {
                const entity = data[i];
                
                // Create entity object for cache
                const entityObj = {
                    id: entity.id,
                    map_id: entity.map_id,
                    x: entity.x,
                    y: entity.y,
                    display_name: entity.display_name,
                    display_org: entity.display_org,
                    search_key: `${entity.display_name} // ${entity.display_org}`.toLowerCase(),
                    marker: null,
                    last_seen: entity.last_seen * 1000,
                    last_seen_human: new Date(entity.last_seen * 1000).toLocaleString('ja-JP', {timeZone: 'Asia/Tokyo'})
                };
                
                // Store in new cache
                newEntityCache[entity.id] = entityObj;
                
                // Skip entities not on current map or out of bounds
                const isOnDifferentMap = entity.map_id !== mapId;
                const isOutOfBounds = entity.x < 0 || entity.y < 0 || entity.x > mapW || entity.y > mapH;
                
                if (isOnDifferentMap || isOutOfBounds) {
                    // Remove marker if it exists
                    if (entityCache[entity.id] && entityCache[entity.id].marker) {
                        map.removeLayer(entityCache[entity.id].marker);
                    }
                    continue;
                }
                
                // Calculate marker position (y-axis is inverted in Leaflet)
                const plotX = entity.x;
                const plotY = mapH - entity.y;
                
                // Create or update marker
                let marker;
                let hadOldMarker = false;
                
                if (entityCache[entity.id] && entityCache[entity.id].marker) {
                    // Update existing marker
                    marker = entityCache[entity.id].marker;
                    hadOldMarker = true;
                } else {
                    // Create new marker
                    const dot = generatePulsatingMarker(15, 'blue');
                    marker = L.marker([plotY, plotX], {icon: dot});
                }
                
                // Create popup content
                const popupHtml = `
                    <div class='card'>
                      <div class='card-user-icon'>
                        <img class='card-user-icon-img' src='./img/user/${entity.id}.png' onerror='this.src="./img/user/user_generic.svg"'>
                      </div>
                      <div class='card-user-text'>
                        <b>${entity.display_name}</b><br />
                        Org: ${entity.display_org}<br />
                        Zone: ${entity.zone_name || 'None'}<br />
                        Last Seen: ${entityObj.last_seen_human}
                      </div>
                    </div>`;
                
                if (hadOldMarker) {
                    // Update existing marker if position changed
                    const positionChanged = entity.x !== entityCache[entity.id].x || entity.y !== entityCache[entity.id].y;
                    if (positionChanged) {
                        marker.slideTo([plotY, plotX], {duration: 500});
                    }
                    
                    // Update tooltip and popup content
                    marker.setTooltipContent(entity.display_name);
                    if (!marker.isPopupOpen()) {
                        marker.setPopupContent(popupHtml);
                    }
                } else {
                    // Configure new marker
                    marker.bindTooltip(entity.display_name, {
                        permanent: true, 
                        direction: 'bottom', 
                        offset: [2, 8]
                    });
                    
                    marker.bindPopup(popupHtml, {
                        maxWidth: '300', 
                        offset: [2, 8]
                    });
                    
                    // Set up event handlers
                    marker.getPopup().on('remove', function() { 
                        markerOpenViaClick = false; 
                    });
                    
                    marker.on('mouseover', function(e) { 
                        if (!markerOpenViaClick) { 
                            this.openPopup(); 
                        } 
                    });
                    
                    marker.on('click', function(e) { 
                        markerOpenViaClick = true; 
                        this.openPopup(); 
                    });
                    
                    marker.on('mouseout', function(e) { 
                        if (!markerOpenViaClick) { 
                            this.closePopup(); 
                        } 
                    });
                    
                    marker.addTo(map);
                }
                
                // Store marker in entity object
                newEntityCache[entity.id].marker = marker;
            }
            
            // Replace entity cache with new data
            Object.keys(entityCache).forEach(key => {
                if (!newEntityCache[key] && entityCache[key].marker) {
                    map.removeLayer(entityCache[key].marker);
                }
            });
            
            Object.assign(entityCache, newEntityCache);
            
            // Update timestamp
            const timeNow = new Date().toLocaleString('ja-JP', {timeZone: 'Asia/Tokyo'});
            const timestampElement = document.getElementById("last-update-time");
            if (timestampElement) {
                timestampElement.textContent = timeNow;
            }
        })
        .fail((jqXHR, textStatus, errorThrown) => {
            console.error("Failed to update entity data:", textStatus, errorThrown);
        });
}

/**
 * Updates zone statistics with current occupancy counts
 */
function updateZone() {
    const zoneApiUrl = `${API_ENDPOINT}/map/${mapId}/zone`;
    
    $.getJSON(zoneApiUrl)
        .done((data) => {
            for (let i = 0; i < data.length; i++) {
                const zone = data[i];
                
                // Skip zones not on current map
                if (zone.map_id !== mapId) {
                    console.warn(`Zone ${zone.name} has different map_id ${zone.map_id} than current map_id ${mapId}`);
                    continue;
                }
                
                // Skip zones not in cache
                if (!(zone.id in zoneCache)) {
                    console.warn(`New zone detected ${zone.id}: ${zone.name} - please reload to see it`);
                    continue;
                }
                
                // Update count if changed
                const oldCount = zoneCache[zone.id];
                const newCount = zone.count;
                
                if (oldCount !== newCount) {
                    const countElement = document.getElementById(zone.id);
                    if (countElement) {
                        countElement.innerHTML = newCount;
                        blink(zone.id);
                        zoneCache[zone.id] = newCount;
                    }
                }
            }
        })
        .fail((jqXHR, textStatus, errorThrown) => {
            console.error("Failed to update zone data:", textStatus, errorThrown);
        });
}


/**
 * Search box and autocomplete functionality
 */
(function() {
    // Search state variables
    let collapseOnBlur = true;
    let activeResult = -1;
    const searchResultDisplayLimit = 10;
    let searchResult = [];
    let searchResultPage = 0;
    let searchKey = "";
    
    // Key codes
    const KEY_ENTER = 13;
    const KEY_LEFT = 37;
    const KEY_UP = 38;
    const KEY_RIGHT = 39;
    const KEY_DOWN = 40;

    /**
     * Initializes the search box functionality
     */
    $.fn.LoadSearchBox = function() {
        $(this).each(function() {
            const element = $(this);
            const appendHtml = `
                <input id="menuButton" class="mdi--menu-open" type="submit" value="" title="Menu">
                <input id="searchBox" class="autocomplete-searchBox" placeholder="Search.."/>
                <input id="searchButton" class="mdi--search" type="submit" value="" title="Search"/>
                <span class="autocomplete-divider"></span>
                <input id="clearButton" class="mdi--clear-bold" type="submit" value="" title="Clear">
            `;
            
            element.addClass("autocomplete-searchContainer");
            element.append(appendHtml);

            // Initialize search box
            $("#searchBox")[0].value = "";
            
            // Set up event handlers
            $("#searchBox").delayKeyup(function(event) {
                switch (event.keyCode) {
                    case KEY_ENTER:
                        searchButtonClick();
                        break;
                    case KEY_UP:
                    case KEY_DOWN:
                    case KEY_LEFT:
                    case KEY_RIGHT:
                        // Do nothing for arrow keys
                        break;
                    default:
                        doSearch();
                        break;
                }
            }, 300);

            $("#searchBox").focus(function() {
                if (isMenuOpen) {
                    clearButtonClick();
                    isMenuOpen = false;
                }

                const resultsDiv = $("#resultsDiv")[0];
                if (resultsDiv !== undefined) {
                    if (activeResult >= 0) {
                        const activeEntity = searchResult[activeResult];
                        if (activeEntity && entityCache[activeEntity]) {
                            const activeMarker = entityCache[activeEntity].marker;
                            if (activeMarker && !activeMarker.isPopupOpen()) {
                                $('#listElement' + activeResult).removeClass('active');
                                activeResult = -1;
                            }
                        }
                    }
                    resultsDiv.style.visibility = "visible";
                } else {
                    doSearch();
                }
            });

            $("#searchBox").blur(function() {
                if (isMenuOpen) {
                    clearButtonClick();
                    isMenuOpen = false;
                }

                const resultsDiv = $("#resultsDiv")[0];
                if (resultsDiv !== undefined) {
                    if (collapseOnBlur) {
                        resultsDiv.style.visibility = "collapse";
                    } else {
                        collapseOnBlur = true;
                        window.setTimeout(function() {
                            $("#searchBox").focus();
                        }, 0);
                    }
                }
            });

            // Button click handlers
            $("#menuButton").click(menuButtonClick);
            $("#searchButton").click(searchButtonClick);
            $("#clearButton").click(clearButtonClick);
        });
    };

    /**
     * Delays keyup event processing to prevent excessive searches
     */
    $.fn.delayKeyup = function(callback, ms) {
        let timer = 0;
        $(this).keyup(function(event) {
            if (event.keyCode !== KEY_ENTER && event.keyCode !== KEY_UP && event.keyCode !== KEY_DOWN) {
                clearTimeout(timer);
                timer = setTimeout(function() {
                    callback(event);
                }, ms);
            } else {
                callback(event);
            }
        });
        return $(this);
    };

    /**
     * Performs search based on current input
     */
    function doSearch() {
        searchKey = $("#searchBox")[0].value.toLowerCase();
        const result = [];
        
        // Find matching entities
        Object.keys(entityCache).forEach(key => {
            const entity = entityCache[key];
            if (searchKey === "" || entity.search_key.includes(searchKey)) {
                result.push(key);
            }
        });

        // Handle no results
        if (result.length < 1) {
            noRecordFoundErr();
            return;
        }

        // Sort results alphabetically
        searchResult = result.sort((a, b) => {
            const keyA = entityCache[a].search_key;
            const keyB = entityCache[b].search_key;
            
            if (keyA < keyB) return -1;
            if (keyA > keyB) return 1;
            return 0;
        });
        
        // Reset result state
        activeResult = -1;
        searchResultPage = 0;

        // Update the results view
        updateResultView();
    }

    /**
     * Updates the search results display
     */
    function updateResultView() {
        const parent = $("#searchBox").parent();

        // Remove existing results and create new container
        $("#resultsDiv").remove();
        parent.append("<div id='resultsDiv' class='autocomplete-result'><ul id='resultList' class='autocomplete-list'></ul><div>");

        // Position the results container
        const resultsDiv = $("#resultsDiv")[0];
        const searchBox = $("#searchBox")[0];
        
        resultsDiv.style.position = searchBox.style.position;
        resultsDiv.style.left = (parseInt(searchBox.style.left) - 10) + "px";
        resultsDiv.style.bottom = searchBox.style.bottom;
        resultsDiv.style.right = searchBox.style.right;
        resultsDiv.style.top = (parseInt(searchBox.style.top) + 25) + "px";
        resultsDiv.style.zIndex = searchBox.style.zIndex;

        // Calculate pagination
        const resultCount = searchResult.length;
        const resultStart = searchResultPage * searchResultDisplayLimit;
        let resultEnd = (searchResultPage + 1) * searchResultDisplayLimit;
        
        if (resultCount < resultEnd) {
            resultEnd = resultCount;
        }
        
        // Add result items
        for (let i = resultStart; i < resultEnd; i++) {
            const entityId = searchResult[i];
            const entry = entityCache[entityId];
            
            if (!entry) continue;
            
            // Determine styling based on entity status
            let contentDetail = `Organization: ${entry.display_org}`;
            let spanClassTitle = "autocomplete-content-text-title";
            let spanClassDetail = "autocomplete-content-text-detail";

            if (entry.marker === null) {
                contentDetail += ` (Last: ${entry.last_seen_human})`;
                spanClassTitle = "autocomplete-content-text-title-offline";
                spanClassDetail = "autocomplete-content-text-detail-offline";
            }
            
            // Create result item HTML
            const html = `
              <li id='listElement${i}' class='autocomplete-listResult'>
                <div id='listElementContent${i}' class='autocomplete-content'>
                  <div class='autocomplete-content-img'>
                    <img src='./img/user/${entry.id}.png' onerror='this.src="./img/user/user_generic.svg"' class='autocomplete-iconStyle' align='middle'>
                  </div>
                  <div class='autocomplete-content-text'>
                    <span class='${spanClassTitle}'>${entry.display_name}</span><br />
                    <span class='${spanClassDetail}'>${contentDetail}</span>
                  </div>
                </div>
              </li>
            `;

            $("#resultList").append(html);

            // Add event handlers to result item
            $("#listElement" + i).mouseenter(function() {
                listElementMouseEnter(this);
            });

            $("#listElement" + i).mouseleave(function() {
                listElementMouseLeave(this);
            });

            $("#listElement" + i).mousedown(function() {
                listElementMouseDown(this);
            });
        }

        // Add pagination info
        const htmlPaging = `
          <div align='right' class='autocomplete-pagingDiv'>
            Showing top ${resultEnd} results out of ${resultCount} results
          </div>
        `;
        $("#resultsDiv").append(htmlPaging);
    }

    /**
     * Event handlers for search result interactions
     */
    function listElementMouseEnter(listElement) {
        const index = parseInt(listElement.id.substr(11));
        if (index !== activeResult) {
            $('#listElement' + index).toggleClass('mouseover');
        }
    }

    function listElementMouseLeave(listElement) {
        const index = parseInt(listElement.id.substr(11));
        if (index !== activeResult) {
            $('#listElement' + index).removeClass('mouseover');
        }
    }

    function listElementMouseDown(listElement) {
        const index = parseInt(listElement.id.substr(11));
        const entityId = searchResult[index];
        const entry = entityCache[entityId];
        
        // Skip if entity has no marker
        if (!entry || entry.marker === null) {
            return;
        }

        if (index !== activeResult) {
            // Update active state
            if (activeResult !== -1) {
                $('#listElement' + activeResult).removeClass('active');
            }

            $('#listElement' + index).removeClass('mouseover');
            $('#listElement' + index).addClass('active');

            activeResult = index;
            entry.marker.openPopup();
        }
    }

    /**
     * Button click handlers
     */
    function clearButtonClick() {
        $("#resultsDiv").remove();
        $("#searchBox")[0].value = "";
        searchResult = [];
        searchResultPage = 0;
        searchKey = "";
        activeResult = -1;
    }

    function searchButtonClick() {
        doSearch();
    }

    function menuButtonClick() {
        if (isMenuOpen) {
            isMenuOpen = false;
            clearButtonClick();
            return;
        }

        // Clear search and create menu
        clearButtonClick();
        const parent = $("#searchBox").parent();

        $("#resultsDiv").remove();
        parent.append("<div id='resultsDiv' class='autocomplete-result'><ul id='resultList' class='autocomplete-list'></ul><div>");

        // Position menu
        const resultsDiv = $("#resultsDiv")[0];
        const searchBox = $("#searchBox")[0];
        
        resultsDiv.style.position = searchBox.style.position;
        resultsDiv.style.left = (parseInt(searchBox.style.left) - 10) + "px";
        resultsDiv.style.bottom = searchBox.style.bottom;
        resultsDiv.style.right = searchBox.style.right;
        resultsDiv.style.top = (parseInt(searchBox.style.top) + 25) + "px";
        resultsDiv.style.zIndex = searchBox.style.zIndex;

        // Add map options to menu
        for (let i = 0; i < maps.length; i++) {
            const html = `
              <li id='listElement${i}' class='autocomplete-listResult'>
                <div id='listElementContent${i}' class='autocomplete-content'>
                  <div class='autocomplete-content-img'>
                      <span class="mdi--map"></span>
                  </div>
                  <div class='autocomplete-content-text'>
                    <span class='autocomplete-content-text-title'>Switch Map: ${maps[i].name}</span>
                  </div>
                </div>
              </li>
            `;
            $("#resultList").append(html);

            // Highlight current map
            if (mapId === maps[i].id) {
                $('#listElement' + i).addClass('active');
            }

            // Add event handlers
            $("#listElement" + i).mouseenter(function() {
                listElementMouseEnter(this);
            });

            $("#listElement" + i).mouseleave(function() {
                listElementMouseLeave(this);
            });

            $("#listElement" + i).mousedown(function() {
                listElementMouseDownMenu(this);
            });
        }

        isMenuOpen = true;
    }

    /**
     * Handles map selection from menu
     */
    function listElementMouseDownMenu(listElement) {
        const index = parseInt(listElement.id.substr(11));
        const entry = maps[index];
        
        if (!entry || !entry.id) {
            return;
        }

        if (mapId !== entry.id) {
            isMenuOpen = false;
            clearButtonClick();
            mapId = entry.id;
            updateMap();
        }
    }

    /**
     * Shows "no results" message
     */
    function noRecordFoundErr() {
        $("#resultsDiv").remove();
        const parent = $("#searchBox").parent();
        const appendHtml = `
          <div id='resultsDiv' class='autocomplete-result'>
            <span class='autocomplete-result-error'>No result for "${searchKey}"</span>
          </div>
        `;

        parent.append(appendHtml);
        activeResult = -1;
        searchResult = [];
        searchResultPage = 0;
        searchKey = "";
    }
    
    /**
     * Pagination functions
     */
    function prevPaging() {
        if (searchResultPage < 1) {
            return;
        }

        $("#searchBox")[0].value = searchKey;
        collapseOnBlur = false;
        activeResult = -1;
        searchResultPage--;
        updateResultView();
    }

    function nextPaging() {
        const resultCount = searchResult.length;
        const nextResultStart = (searchResultPage + 1) * searchResultDisplayLimit;
        
        if (resultCount < nextResultStart) {
            return;
        }

        $("#searchBox")[0].value = searchKey;
        collapseOnBlur = false;
        activeResult = -1;
        searchResultPage++;
        updateResultView();
    }
})();
