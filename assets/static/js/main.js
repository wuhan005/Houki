async function startProxy(address) {
    await fetch('/proxy/start', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({Address: address})
    });
    location.reload();
}

async function stopProxy() {
    await fetch('/proxy/shut-down', {
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
