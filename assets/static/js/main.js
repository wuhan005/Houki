async function startProxy(address) {
    await fetch('/proxy/start', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({Address: address})
    }).then(async res => {
        if (res.status === 200) {
            location.reload();
            return
        }
        alert((await res.json()).msg)
    }).catch(err => {
        console.log(err)
    })
}

async function stopProxy() {
    await fetch('/proxy/shut-down', {
        method: 'POST',
    })
    location.reload();
}

async function openBrowser() {
    await fetch('/chromedp/', {
        method: 'POST',
    })
}

async function newModule() {
    let moduleID = document.getElementById('module-id').value;
    let body = {}
    try {
        body = JSON.parse(editor.getValue())
    } catch (e) {
        alert('Invalid JSON')
        return
    }

    await fetch('/modules/new', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ID: moduleID, Body: body})
    }).then(async res => {
        if (res.status === 200) {
            window.location.href = '/modules/'
            return
        }
        alert((await res.json()).msg)
    }).catch(err => {
        console.log(err)
    })
}

function updateModule(id) {
    let body = {}
    try {
        body = JSON.parse(editor.getValue())
    } catch (e) {
        alert('Invalid JSON')
        return
    }

    fetch(`/modules/${id}`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ID: id, Body: body})
    }).then(async res => {
        if (res.status === 200) {
            location.reload();
            return
        }
        alert((await res.json()).msg)
    }).catch(err => {
        console.log(err)
    })
}

async function deleteModule(id) {
    await fetch(`/modules/${id}`, {
        method: 'DELETE',
    }).then(async res => {
        if (res.status === 200) {
            window.location.href = '/modules/'
            return
        }
        alert((await res.json()).msg)
    }).catch(err => {
        console.log(err)
    })
}

function enableModule(id) {
    fetch(`/modules/${id}/enable`, {
        method: 'POST',
    }).then(async res => {
        if (res.status === 200) {
            location.reload();
            return
        }
        alert((await res.json()).msg)
    }).catch(err => {
        console.log(err)
    })
}

function disableModule(id) {
    fetch(`/modules/${id}/disable`, {
        method: 'POST',
    }).then(async res => {
        if (res.status === 200) {
            location.reload();
            return
        }
        alert((await res.json()).msg)
    }).catch(err => {
        console.log(err)
    })
}
