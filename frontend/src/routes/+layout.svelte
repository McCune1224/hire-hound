<script>
	import '../app.css';
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<nav />
<h1>FOOBARBAZX</h1>

<slot />

<style lang="postcss">
	:global(html) {
		background-color: theme(colors.zinc.900);
		color: theme(colors.white);
	}
</style>
