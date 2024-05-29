module.exports = async ({
  core
}) => {
  const fs = require('fs');

  // Read the text from the file
  fs.readFile('coverage/integration.report', 'utf8', (err, data) => {
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

    const report = jsonData.filter(data => {
      const driftDetectedSubStrings = [
        'the plan was not empty',
        'expected an error but got none',
      ];

      if (!data.Output) {
        return false;
      }

      let driftDetected = false
      for (let i = 0; i < driftDetectedSubStrings.length; i++) {
        if (data.Output.includes(driftDetectedSubStrings[i])) {
          return true;
        }
      }

      return driftDetected
    }).map(t => t.Output.trim());

    let msg = 'No drift detected.';
    if (report.length > 0) {
      msg = `'${report.join('\n\n')}'`;
    }

    console.log('drift_report:', msg);

    core.setOutput('drift_report', msg);
  });
};
