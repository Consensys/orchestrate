const { Octokit } = require("@octokit/rest");

if (process.argv.length != 4) {
  console.error('Expected 2 arguments: PR Title and source branch.');
  process.exit(1);
}

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

console.log(process.argv);

octokit.pulls.create({
  owner:'ConsenSys',
  repo:'doc.orchestrate',
  title:process.argv[2],
  head:process.argv[3],
  base:'master',
  maintainer_can_modify:true,
  draft:true
});
