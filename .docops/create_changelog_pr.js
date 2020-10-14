const { Octokit } = require("@octokit/rest");

if (process.argv.length != 4) {
  console.error('Expected 2 arguments: PR Title and source branch.');
  process.exit(1);
}

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

console.log(process.argv);
console.log(octokit);

// octokit.pulls.create({
//   'ConsenSys',
//   'doc.orchestrate',
//   process.argv[2],
//   process.argv[3],
//   'master',
// });
