module.exports = async ({
  core
}) => {
  const fs = require('fs');

  // Read the text from the file
  fs.readFile('coverage/unit.report', 'utf8', (err, data) => {
    if (err) {
      console.error('Error reading file:', err);
      return;
    }

    // Split the text into individual lines
    const lines = data.trim().split('\n');

    // Parse each line and convert it to JSON
    const jsonData = lines.map(line => {
      try {
        return JSON.parse(line);
      } catch (error) {
        console.error('Error parsing line:', error);
        return null;
      }
    });

    // Print the resulting JSON array
    // console.log(jsonData);

    const report = jsonData.filter(data => {
      return data.Output ? data.Output.includes('error: After applying this test step, the plan was not empty') : false;
    }).map(t => t.Output.trim());

    console.log(report);

    let msg = 'error: After applying this test step, the plan was not empty';
    if (report.length > 0) {
      msg = `'${report.join('\n')}'`;
    }

    core.setOutput('failed_tests_with_drift', msg);

    // // Convert JSON data to string
    // const jsonString = JSON.stringify(failedTests, null, 2);

    // // Write the JSON data to a file
    // fs.writeFile('output.json', jsonString, 'utf8', err => {
    //   if (err) {
    //     console.error('Error writing file:', err);
    //     return;
    //   }
    //   console.log('JSON data has been written to output.json');
    // });
  });
};
