const reporter = require('cucumber-html-reporter');

const metadata = Object.keys(process.env)
    .filter(key => key.split("_")[0] === "METADATA" && process.env[key] !== "")
    .reduce((obj, key) => {
        const name = key.split("_")[1]
        obj[name] = process.env[key];
        return obj;
    }, { alias: process.env.CUCUMBER_ALIAS });

const options = {
    brandTitle: 'PegaSys Orchestrate end-to-end tests',
    theme: 'bootstrap',
    jsonFile: process.env.CUCUMBER_INPUT ? process.env.CUCUMBER_INPUT : 'report.json',
    output: process.env.CUCUMBER_OUTPUT ? process.env.CUCUMBER_OUTPUT : 'report.html',
    reportSuiteAsScenarios: true,
    metadata
};

reporter.generate(options);
