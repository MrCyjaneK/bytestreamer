function pteLogin(username, password) {
    if (username == undefined || password == undefined) {
        username = document.getElementById('username').value
        password = document.getElementById('password').value
    }
    localStorage.username = document.getElementById("username").value
    localStorage.password = document.getElementById("password").value
    localStorage.apikey = document.getElementById("apikey").value
    fetch('https://pte.nu/apibslogin', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'email='+encodeURIComponent(username)+'&password='+encodeURIComponent(password)
    }).then(r => r.json())
    .then(r => {
        if (r.error_code) {
            document.getElementById('error').innerText = r.message
            return
        }
        console.log(r)
        localStorage.kkey = r.key
        localStorage.lang = r.lang
        localStorage.userid = r.userid
        localStorage.apikey = document.getElementById('apikey').value
        window.location.href = "/player.html"
    })
}

function pteSearch(query) {
    if (query == undefined) {
        query = document.getElementById("search").value
    }
    fetch("https://pte.nu/apistreamtorrents?id="+encodeURIComponent(localStorage.userid)+"&key="+encodeURIComponent(localStorage.kkey)+"&category=1&page=0&search="+encodeURIComponent(query))
    .then(r => r.json())
    .then(r => {
        let elm = document.getElementById('stuff')
        elm.innerHTML = ""
        console.log(r)
        if (!r.torrents) {
            elm.innerText = "No torrents found! Oupsie!"
            return
        }
        let ul = document.createElement('ul')
        for (i in r.torrents) {
            let li = document.createElement('li')
            let a = document.createElement('a')
            console.log(r.torrents[i])
            a.href = "/info/pte.html?id="+encodeURIComponent(r.torrents[i].id)+"&title="+encodeURIComponent(r.torrents[i].real_name)
            a.innerText = r.torrents[i].name
            li.appendChild(a)
            ul.appendChild(li)
        }
        elm.appendChild(ul)
        fetch("https://pte.nu/apistreamtorrents?id="+encodeURIComponent(localStorage.userid)+"&key="+encodeURIComponent(localStorage.kkey)+"&category=2&page=0&search="+encodeURIComponent(query))
    .then(r => r.json())
    .then(r => {
        let elm = document.getElementById('stuff')
        console.log(r)
        if (!r.torrents) {
            elm.innerText = "No torrents found! Oupsie!"
            return
        }
        elm.appendChild(document.createElement('hr'))
        let ul = document.createElement('ul')
        for (i in r.torrents) {
            let li = document.createElement('li')
            let a = document.createElement('a')
            console.log(r.torrents[i])
            a.href = "/info/pte.html?id="+encodeURIComponent(r.torrents[i].id)+"&title="+encodeURIComponent(r.torrents[i].real_name)
            a.innerText = r.torrents[i].name
            li.appendChild(a)
            ul.appendChild(li)
        }
        elm.appendChild(ul)
    })
    })
}

function loadInfoPte() {
    let id = getParameterByName('id')
    let title = getParameterByName('title')
    document.getElementById('title').innerText = title
    fetch("/api/torrent_info?torrentid="+encodeURIComponent(id)+"&apikey="+encodeURIComponent(localStorage.apikey))
    .then(r => r.json())
    .then(r => {
        let ul = document.createElement('ul')
        for (i in r.Files) {
            let li = document.createElement('li')
            let a = document.createElement('a')
            a.target = "_blank"
            a.href = "/play?torrentid="+encodeURIComponent(id)+"&file="+encodeURIComponent(r.Files[i].Path.join("/"))
            a.innerText = r.Files[i].Path[r.Files[i].Path.length-1]
            li.appendChild(a)
            ul.appendChild(li)
        }
        document.getElementById("files").appendChild(ul)
    })
    fetch("https://pte.nu/bsgettorrent/"+encodeURIComponent(id)+"?id="+encodeURIComponent(localStorage.userid)+"&key="+encodeURIComponent(localStorage.kkey))
    .then(r => r.json())
    .then(r => {
        console.log(r)
        let ul = document.createElement('ul')
        for (i in r.actors_name) {
            let li = document.createElement('li')
            li.innerText = r.actors_name[i]
            ul.appendChild(li)
        }
        document.getElementById('actors').appendChild(ul)
        document.getElementById('description').innerText = r.description_pl || r.description
        let torrents = JSON.parse("["+r.others.substr(1,r.others.length - 2)+"]")
        ul = document.createElement('ul')
        for (i in torrents) {
            let t = torrents[i].substr(1, torrents[i].length - 2).split(',')
            console.log(t)
            let li = document.createElement('li')
            let a = document.createElement('a')
            a.href = "/info/pte.html?id="+encodeURIComponent(t[3])+"&title="+encodeURIComponent(t[0])
            a.innerText = t[0]+" | s: "+t[1]+" | "+humanFileSize(t[2])
            li.appendChild(a)
            ul.appendChild(li)
        }
        document.getElementById('torrents').appendChild(ul)
    })
}

function getParameterByName(name) {
    var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
    return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
}

function humanFileSize(bytes, si=false, dp=1) {
    const thresh = si ? 1000 : 1024;
  
    if (Math.abs(bytes) < thresh) {
      return bytes + ' B';
    }
  
    const units = si 
      ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'] 
      : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
    let u = -1;
    const r = 10**dp;
  
    do {
      bytes /= thresh;
      ++u;
    } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);
  
  
    return bytes.toFixed(dp) + ' ' + units[u];
  }
  