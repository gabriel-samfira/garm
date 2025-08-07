<script lang="ts">
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { garmApi } from '$lib/api/client.js';
	import type { Repository, Organization, Pool, Instance } from '$lib/api/types.js';

	let loading = true;
	let stats = {
		repositories: 0,
		organizations: 0,
		pools: 0,
		instances: 0
	};
	let error = '';

	onMount(async () => {
		try {
			const [repos, orgs, pools, instances] = await Promise.all([
				garmApi.listRepositories(),
				garmApi.listOrganizations(),
				garmApi.listPools(),
				garmApi.listInstances()
			]);

			stats = {
				repositories: repos.length,
				organizations: orgs.length,
				pools: pools.length,
				instances: instances.length
			};
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load dashboard data';
			console.error('Dashboard error:', err);
		} finally {
			loading = false;
		}
	});

	$: statCards = [
		{
			title: 'Repositories',
			value: stats.repositories,
			icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
			color: 'blue',
			href: `${base}/repositories`
		},
		{
			title: 'Organizations', 
			value: stats.organizations,
			icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z',
			color: 'green',
			href: `${base}/organizations`
		},
		{
			title: 'Pools',
			value: stats.pools,
			icon: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z',
			color: 'purple',
			href: `${base}/pools`
		},
		{
			title: 'Instances',
			value: stats.instances,
			icon: 'M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z',
			color: 'yellow',
			href: `${base}/instances`
		}
	];

	function getColorClasses(color: string) {
		const colorMap = {
			blue: 'bg-blue-500 text-white',
			green: 'bg-green-500 text-white',
			purple: 'bg-purple-500 text-white',
			yellow: 'bg-yellow-500 text-white'
		};
		return colorMap[color as keyof typeof colorMap] || 'bg-gray-500 text-white';
	}
</script>

<svelte:head>
	<title>Dashboard - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div>
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
		<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
			Welcome to GARM - GitHub Actions Runner Manager
		</p>
	</div>

	{#if loading}
		<!-- Loading state -->
		<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
			{#each Array(4) as _}
				<div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg animate-pulse">
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 bg-gray-300 dark:bg-gray-600 rounded"></div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-24 mb-2"></div>
								<div class="h-8 bg-gray-300 dark:bg-gray-600 rounded w-16"></div>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{:else if error}
		<!-- Error state -->
		<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
					</svg>
				</div>
				<div class="ml-3">
					<h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading dashboard</h3>
					<p class="mt-2 text-sm text-red-700 dark:text-red-300">{error}</p>
				</div>
			</div>
		</div>
	{:else}
		<!-- Stats cards -->
		<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
			{#each statCards as card}
				<a
					href={card.href}
					class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow duration-200"
				>
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 rounded-md {getColorClasses(card.color)} flex items-center justify-center">
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={card.icon}></path>
									</svg>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
										{card.title}
									</dt>
									<dd class="text-lg font-medium text-gray-900 dark:text-white">
										{card.value}
									</dd>
								</dl>
							</div>
						</div>
					</div>
				</a>
			{/each}
		</div>

		<!-- Quick actions -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="p-6">
				<h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">Quick Actions</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
					Common tasks you can perform
				</p>
				
				<div class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
					<a
						href="{base}/repositories"
						class="relative block w-full bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg p-6 hover:border-gray-400 dark:hover:border-gray-500 hover:shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
					>
						<div class="flex items-center">
							<svg class="flex-shrink-0 w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
							</svg>
							<div class="ml-4">
								<h4 class="text-sm font-medium text-gray-900 dark:text-white">Add Repository</h4>
								<p class="text-sm text-gray-500 dark:text-gray-400">Connect a new GitHub repository</p>
							</div>
						</div>
					</a>

					<a
						href="{base}/pools"
						class="relative block w-full bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg p-6 hover:border-gray-400 dark:hover:border-gray-500 hover:shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
					>
						<div class="flex items-center">
							<svg class="flex-shrink-0 w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
							</svg>
							<div class="ml-4">
								<h4 class="text-sm font-medium text-gray-900 dark:text-white">Create Pool</h4>
								<p class="text-sm text-gray-500 dark:text-gray-400">Set up a new runner pool</p>
							</div>
						</div>
					</a>

					<a
						href="{base}/instances"
						class="relative block w-full bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg p-6 hover:border-gray-400 dark:hover:border-gray-500 hover:shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
					>
						<div class="flex items-center">
							<svg class="flex-shrink-0 w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
							</svg>
							<div class="ml-4">
								<h4 class="text-sm font-medium text-gray-900 dark:text-white">View Instances</h4>
								<p class="text-sm text-gray-500 dark:text-gray-400">Monitor running instances</p>
							</div>
						</div>
					</a>
				</div>
			</div>
		</div>
	{/if}
</div>