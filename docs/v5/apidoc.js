'use strict';

window.onload = function () {
    registerMethodFilter();
    registerTagFilter();
    registerLanguageFilter();
    initExample()
};

function registerMethodFilter() {
    const menu = document.querySelector('.methods-selector');

    menu.style.display = 'block';

    const list = menu.querySelectorAll('li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const method = event.target.parentNode.parentNode.getAttribute('data-method');

            const apis = document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (api.getAttribute('data-method') !== method) {
                    return;
                }

                api.style.display = chk ? 'block' : 'none';
            });
        });
    });
}

function registerTagFilter() {
    const menu = document.querySelector('.tags-selector');

    menu.style.display = 'block';

    const list = menu.querySelectorAll('li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-tag');

            const apis = document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (!api.getAttribute('data-tags').includes(tag + ',')) {
                    return;
                }

                api.style.display = chk ? 'block' : 'none';
            });
        });
    });
}

function registerLanguageFilter() {
    const menu = document.querySelector('.languages-selector');

    menu.style.display = 'block';

    const list = menu.querySelectorAll('li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            if (!event.target.checked) {
                return;
            }

            const lang = event.target.parentNode.parentNode.getAttribute('lang');
            const elems = document.querySelectorAll('[data-locale]');

            elems.forEach((elem) => {
                if (elem.getAttribute('lang') === lang) {
                    elem.className = '';
                } else {
                    elem.className = 'hidden';
                }
            });
        });
    });
}

function initExample() {
    const buttons = document.querySelectorAll('.toggle-example');

    buttons.forEach((btn)=>{
        btn.addEventListener('click', (e)=> {
            const parent = e.target.parentNode.parentNode.parentNode;
            const table = parent.querySelector('.param-list');
            const pre = parent.querySelector('.example');

            if (table.getAttribute('data-visible') === 'true') {
                table.setAttribute('data-visible', 'false');
                table.style.display = 'none';
            } else {
                table.setAttribute('data-visible', 'true');
                table.style.display = 'table';
            }

            if (pre.getAttribute('data-visible') === 'true') {
                pre.setAttribute('data-visible', 'false');
                pre.style.display = 'none';
            } else {
                pre.setAttribute('data-visible', 'true');
                pre.style.display = 'block';
            }
        });
    });
}
