var nodeenvconfiguration = require('node-env-configuration');
var reporter = require('cucumber-html-reporter');

var defaults = {
    theme: 'bootstrap',
    jsonFile: 'in/report.json',
    output: 'out/report.html',
    reportSuiteAsScenarios: true,
};

// see https://github.com/whynotsoluciones/node-env-configuration 
var options = nodeenvconfiguration({
    defaults: defaults,
    prefix: 'report'
});

reporter.generate(options);
