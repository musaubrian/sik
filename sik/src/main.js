const { invoke } = window.__TAURI__.core

let searchInputEl
let searchResultEl

async function openUrl(url) {
  if (url.length > 1) {
    await invoke('open_file', {
      url: url,
    })
  }
}

async function searchQuery() {
  searchResultEl.innerHTML = ''

  const result = await invoke('search', {
    query: searchInputEl.value,
  })
  const jsonRes = JSON.parse(result)
  if (jsonRes.length < 1) {
    const noResultSpan = document.createElement('span')
    noResultSpan.id = 'noResult'
    noResultSpan.innerText = `:( Nothing found matching <${searchInputEl.value}>`

    searchResultEl.appendChild(noResultSpan)
  }
  for (let idx = 0; idx < jsonRes.length; idx++) {
    const element = jsonRes[idx]

    const resultDiv = document.createElement('div')
    resultDiv.id = 'single-result'

    const locationSpan = document.createElement('span')
    locationSpan.className = 'location'
    locationSpan.innerText = element.full_path

    const fileLink = document.createElement('a')
    fileLink.href = `file://${element.full_path}`
    fileLink.innerText = element.filename

    // Add click event listener to the link
    fileLink.addEventListener('click', function (e) {
      e.preventDefault()
      const hrefValue = this.getAttribute('href')
      openUrl(hrefValue)
    })

    const hitsList = document.createElement('ul')
    hitsList.className = 'hits-list'

    element.hits.forEach((hit) => {
      const hitItem = document.createElement('li')
      hitItem.className = 'hit-item'
      hitItem.innerHTML = hit.trim()
      hitsList.appendChild(hitItem)
    })

    resultDiv.appendChild(locationSpan)
    resultDiv.appendChild(fileLink)
    resultDiv.appendChild(hitsList)

    searchResultEl.appendChild(resultDiv)
  }
}

window.addEventListener('DOMContentLoaded', () => {
  searchInputEl = document.querySelector('#search-input')
  searchResultEl = document.querySelector('#search-results')
  document.querySelector('#search-form').addEventListener('submit', (e) => {
    e.preventDefault()
    searchQuery()
  })
})
