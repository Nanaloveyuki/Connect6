import { spawn } from 'node:child_process';

const [, , task, ...rawTargets] = process.argv;

const allowedTasks = new Set(['dev', 'check', 'build']);
const allowedTargets = ['backend', 'frontend'];
const runner =
	process.platform === 'win32'
		? { command: 'cmd.exe', prefix: ['/d', '/s', '/c', 'pnpm'] }
		: { command: 'pnpm', prefix: [] };

if (!allowedTasks.has(task)) {
	console.error(`Unsupported task: ${task}`);
	process.exit(1);
}

const targets = resolveTargets(rawTargets);

if (targets.length === 1) {
	const exitCode = await runCommand(targets[0], task);
	process.exit(exitCode);
}

if (task === 'dev') {
	const exitCode = await runParallel(targets, task);
	process.exit(exitCode);
}

for (const target of targets) {
	const exitCode = await runCommand(target, task);
	if (exitCode !== 0) {
		process.exit(exitCode);
	}
}

function resolveTargets(values) {
	if (values.length === 0) {
		return allowedTargets;
	}

	const normalized = values.map((value) => value.toLowerCase());
	for (const target of normalized) {
		if (!allowedTargets.includes(target)) {
			console.error(`Unsupported target: ${target}`);
			process.exit(1);
		}
	}

	return [...new Set(normalized)];
}

function runCommand(target, taskName) {
	return new Promise((resolve) => {
		const child = spawn(runner.command, [...runner.prefix, '--dir', target, taskName], {
			cwd: process.cwd(),
			stdio: 'inherit'
		});

		child.on('exit', (code, signal) => {
			if (signal) {
				resolve(1);
				return;
			}

			resolve(code ?? 1);
		});
	});
}

function runParallel(targetsList, taskName) {
	return new Promise((resolve) => {
		const children = targetsList.map((target) => {
			return spawn(runner.command, [...runner.prefix, '--dir', target, taskName], {
				cwd: process.cwd(),
				stdio: 'inherit'
			});
		});

		let settled = false;
		let closedCount = 0;

		const terminateOthers = (sourceChild) => {
			for (const child of children) {
				if (child !== sourceChild && !child.killed) {
					child.kill('SIGTERM');
				}
			}
		};

		const finish = (code) => {
			if (!settled) {
				settled = true;
				resolve(code);
			}
		};

		for (const signal of ['SIGINT', 'SIGTERM']) {
			process.on(signal, () => {
				terminateOthers(undefined);
				finish(1);
			});
		}

		children.forEach((child) => {
			child.on('exit', (code, signal) => {
				closedCount += 1;

				if (!settled && (signal || code !== 0)) {
					terminateOthers(child);
					finish(code ?? 1);
					return;
				}

				if (closedCount === children.length) {
					finish(0);
				}
			});
		});
	});
}
