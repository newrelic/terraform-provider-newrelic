var assert = require('assert');

$browser.get("https://staging-login.newrelic.com/").then(function(){
    $browser.findElement($driver.By.name("login[email]")).sendKeys("skumarreddy@newrelic.com").then(function(){
        $browser.findElement($driver.By.name("login[password]")).sendKeys($secure.TF_PASSWERTY).then(function(){
            $browser.findElement($driver.By.name("button")).click();
        });
    });
});
