<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<div class="app">
	<header class="topbar">
		<a class="brand" href="/">
			<span class="brand-mark">✉</span>
			<span class="brand-name">Box</span>
		</a>
		<nav class="nav">
			<a href="/">Inbox</a>
		</nav>
	</header>
	<main class="content">
		{@render children?.()}
	</main>
</div>
