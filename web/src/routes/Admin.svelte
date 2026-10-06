<script lang="ts">
  import { api } from '../lib/api';
  import { CATS, type Content, type Work } from '../lib/types';

  type Tab = 'profile' | 'experience' | 'tools' | 'works' | 'contacts';
  const tabs: [Tab, string][] = [
    ['profile', 'Профиль'], ['experience', 'Опыт'], ['tools', 'Инструменты'], ['works', 'Работы'], ['contacts', 'Контакты'],
  ];

  let authed = $state<boolean | null>(null);
  let password = $state('');
  let loginErr = $state('');
  let tab = $state<Tab>('profile');
  let c = $state<Content | null>(null);
  let msg = $state('');
  let busy = $state(false);
  let editing = $state<Work | null>(null);

  const emptyWork = (): Work => ({ id: '', cat: 'home', type: 'Квартира', title: '', area: '', year: String(new Date().getFullYear()), scope: '', cover: '', pdf: '', size: '' });

  async function load() {
    c = await api.content();
  }
  api.session().then(async (s) => {
    authed = s.authed;
    if (s.authed) await load();
  });

  function flash(text: string) {
    msg = text;
    setTimeout(() => (msg === text ? (msg = '') : 0), 3000);
  }

  async function run(fn: () => Promise<void>, ok = 'Сохранено') {
    busy = true;
    try { await fn(); flash(ok); } catch (e) { flash('Ошибка: ' + (e as Error).message); }
    busy = false;
  }

  async function login(e: SubmitEvent) {
    e.preventDefault();
    loginErr = '';
    try {
      await api.login(password);
      password = '';
      authed = true;
      await load();
    } catch (err) { loginErr = (err as Error).message; }
  }

  async function logout() {
    await api.logout();
    authed = false;
    c = null;
  }

  const save = () => run(async () => { c = await api.saveContent(c!); });

  async function pick(e: Event, apply: (url: string, size: string) => void) {
    const input = e.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    await run(async () => { const r = await api.upload(file); apply(r.url, r.size); }, 'Файл загружен (не забудьте сохранить)');
    input.value = '';
  }

  async function saveWork() {
    const w = editing!;
    await run(async () => {
      if (w.id) await api.updateWork(w); else await api.createWork(w);
      await load();
      editing = null;
    });
  }

  async function removeWork(w: Work) {
    if (!confirm(`Удалить «${w.title}»?`)) return;
    await run(async () => { await api.deleteWork(w.id); await load(); }, 'Удалено');
  }

  async function move(i: number, d: number) {
    const list = [...c!.works];
    const j = i + d;
    if (j < 0 || j >= list.length) return;
    [list[i], list[j]] = [list[j], list[i]];
    c!.works = list;
    await run(() => api.orderWorks(list.map((w) => w.id)), 'Порядок сохранён');
  }
</script>

<div class="admin">
  {#if authed === null}
    <p>Загрузка…</p>
  {:else if !authed}
    <form class="login" onsubmit={login}>
      <h1>Админка</h1>
      <label>Пароль<input type="password" bind:value={password} autocomplete="current-password" required /></label>
      {#if loginErr}<div class="err">{loginErr}</div>{/if}
      <button class="primary">Войти</button>
    </form>
  {:else if c}
    <header>
      <strong>Админка · ФЖ</strong>
      <a href="/" target="_blank">Открыть сайт ↗</a>
      <span class="msg" role="status">{msg}</span>
      <button onclick={logout}>Выйти</button>
    </header>
    <nav>
      {#each tabs as [k, label]}<button class:on={tab === k} onclick={() => (tab = k)}>{label}</button>{/each}
    </nav>

    <main>
      {#if tab === 'profile'}
        <label>Имя<input bind:value={c.profile.name} /></label>
        <label>Описание<textarea rows="3" bind:value={c.profile.tagline}></textarea></label>
        <label class="check"><input type="checkbox" bind:checked={c.profile.available} /> Беру новые объекты</label>
        <fieldset><legend>Статистика</legend>
          {#each c.profile.stats as s}
            <div class="row2"><input placeholder="Подпись" bind:value={s.label} /><input placeholder="Значение" bind:value={s.value} /></div>
          {/each}
        </fieldset>
        <fieldset><legend>Фото на главной</legend>
          {#if c.profile.heroImage}<img class="thumb" src={c.profile.heroImage} alt="" />{/if}
          <input type="file" accept="image/jpeg,image/png,image/webp" onchange={(e) => pick(e, (u) => (c!.profile.heroImage = u))} />
        </fieldset>
        <fieldset><legend>Резюме (PDF)</legend>
          {#if c.profile.resumePdf}<a href={c.profile.resumePdf} target="_blank">Текущий файл</a>
            <button onclick={() => (c!.profile.resumePdf = '')}>Убрать</button>{/if}
          <input type="file" accept="application/pdf" onchange={(e) => pick(e, (u) => (c!.profile.resumePdf = u))} />
        </fieldset>
      {:else if tab === 'experience'}
        {#each c.experience as e, i}
          <fieldset>
            <div class="row2"><input placeholder="Период" bind:value={e.period} /><input placeholder="Должность" bind:value={e.role} /></div>
            <input placeholder="Компания" bind:value={e.company} />
            <textarea rows="3" placeholder="Описание" bind:value={e.desc}></textarea>
            <button onclick={() => c!.experience.splice(i, 1)}>Удалить</button>
          </fieldset>
        {/each}
        <button onclick={() => c!.experience.push({ period: '', role: '', company: '', desc: '' })}>+ Добавить место</button>
      {:else if tab === 'tools'}
        {#each c.tools as t, i}
          <div class="row3">
            <input placeholder="Название" bind:value={t.name} />
            <input type="number" min="1" max="5" bind:value={t.level} />
            <button onclick={() => c!.tools.splice(i, 1)}>×</button>
          </div>
        {/each}
        <button onclick={() => c!.tools.push({ name: '', level: 3 })}>+ Инструмент</button>
        <label>Навыки (по одному в строке)
          <textarea rows="10" value={c.skills.join('\n')} oninput={(e) => (c!.skills = e.currentTarget.value.split('\n').map((s) => s.trim()).filter(Boolean))}></textarea>
        </label>
      {:else if tab === 'contacts'}
        <label>Email<input type="email" bind:value={c.contacts.email} /></label>
        <label>Telegram (ссылка)<input bind:value={c.contacts.telegram} placeholder="https://t.me/…" /></label>
        <label>Телефон<input bind:value={c.contacts.phone} /></label>
      {:else if tab === 'works'}
        {#if editing}
          <fieldset>
            <legend>{editing.id ? 'Редактирование' : 'Новая работа'}</legend>
            <label>Название<input bind:value={editing.title} /></label>
            <div class="row2">
              <label>Категория
                <select bind:value={editing.cat}>{#each CATS.slice(1) as k}<option value={k.key}>{k.label}</option>{/each}</select>
              </label>
              <label>Тип (подпись)<input bind:value={editing.type} /></label>
            </div>
            <div class="row2"><label>Площадь<input bind:value={editing.area} placeholder="142 м²" /></label><label>Год<input bind:value={editing.year} /></label></div>
            <label>Описание<textarea rows="3" bind:value={editing.scope}></textarea></label>
            <div>Обложка:
              {#if editing.cover}<img class="thumb" src={editing.cover} alt="" />{/if}
              <input type="file" accept="image/jpeg,image/png,image/webp" onchange={(e) => pick(e, (u) => (editing!.cover = u))} />
            </div>
            <div>PDF проекта:
              {#if editing.pdf}<a href={editing.pdf} target="_blank">файл</a> ({editing.size}){/if}
              <input type="file" accept="application/pdf" onchange={(e) => pick(e, (u, s) => { editing!.pdf = u; editing!.size = s; })} />
            </div>
            <div class="actions">
              <button class="primary" disabled={busy} onclick={saveWork}>Сохранить работу</button>
              <button onclick={() => (editing = null)}>Отмена</button>
            </div>
          </fieldset>
        {:else}
          <button class="primary" onclick={() => (editing = emptyWork())}>+ Новая работа</button>
          {#each c.works as w, i (w.id)}
            <div class="work">
              {#if w.cover}<img class="thumb" src={w.cover} alt="" />{:else}<div class="thumb"></div>{/if}
              <div class="grow"><strong>{w.title}</strong><br /><small>{w.type} · {w.area} · {w.year}{w.pdf ? ' · PDF ' + w.size : ''}</small></div>
              <button onclick={() => move(i, -1)} aria-label="Выше">↑</button>
              <button onclick={() => move(i, 1)} aria-label="Ниже">↓</button>
              <button onclick={() => (editing = { ...w })}>Изменить</button>
              <button onclick={() => removeWork(w)}>Удалить</button>
            </div>
          {/each}
        {/if}
      {/if}
    </main>

    {#if tab !== 'works'}
      <footer><button class="primary" disabled={busy} onclick={save}>Сохранить</button></footer>
    {/if}
  {/if}
</div>

<style>
  .admin { max-width: 860px; margin: 0 auto; padding: 24px 16px 96px; font: 15px/1.4 var(--sans); color: var(--ink); }
  header { display: flex; gap: 16px; align-items: center; padding-bottom: 12px; border-bottom: 1px solid var(--ink); }
  .msg { flex: 1; text-align: right; color: var(--accent); }
  nav { display: flex; flex-wrap: wrap; gap: 4px; margin: 16px 0; }
  nav button.on { background: var(--ink); color: var(--paper); }
  main { display: flex; flex-direction: column; gap: 14px; }
  label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
  label.check { flex-direction: row; align-items: center; gap: 8px; font-size: 15px; }
  input, textarea, select { font: inherit; padding: 8px 10px; border: 1px solid var(--ink); background: var(--paper); color: var(--ink); border-radius: 0; width: 100%; }
  input[type='checkbox'] { width: auto; }
  input[type='file'] { border: 0; padding: 8px 0; }
  fieldset { border: 1px solid var(--ink); display: flex; flex-direction: column; gap: 10px; padding: 14px; margin: 0; }
  button { font: inherit; padding: 8px 14px; border: 1px solid var(--ink); background: var(--paper); color: var(--ink); cursor: pointer; }
  button:hover:not(:disabled) { background: var(--hover); }
  button.primary { background: var(--ink); color: var(--paper); }
  button:disabled { opacity: 0.5; }
  .row2 { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
  .row3 { display: grid; grid-template-columns: 1fr 80px 40px; gap: 10px; }
  .work { display: flex; gap: 8px; align-items: center; border: 1px solid var(--ink); padding: 8px; flex-wrap: wrap; }
  .grow { flex: 1; min-width: 160px; }
  .thumb { width: 72px; height: 72px; object-fit: cover; border: 1px solid var(--ink); background: var(--hover); }
  .actions { display: flex; gap: 8px; }
  footer { position: fixed; left: 0; right: 0; bottom: 0; padding: 12px 16px; background: var(--paper); border-top: 1px solid var(--ink); display: flex; justify-content: center; }
  .login { max-width: 320px; margin: 15vh auto 0; display: flex; flex-direction: column; gap: 14px; }
  .err { color: #a33; }
  @media (max-width: 600px) { .row2 { grid-template-columns: 1fr; } }
</style>
