<?php

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";

session_start();
$sessionId = session_id();

$command = __DIR__ . '/../bin/nli';
$configPath = __DIR__ . '/../resources/dbpedia';
$varDir = __DIR__ . '/../var';
$fullCommand = sprintf('%s -s %s -c %s -d %s -r json "%s"', $command, $sessionId, $configPath, $varDir, $query);
$start = microtime(true);

$process = exec($fullCommand, $output);

$end = microtime(true);
$duration = sprintf("%.2f", $end - $start);
$result = json_decode(implode($output), true);

$handle = fopen(__DIR__ . '/../var/log/' . date('Y-m') . '-queries.log', 'a');

$answer = $result['Answer'];
$optionKeys = $result['OptionKeys'];
$optionValues = $result['OptionValues'];
$errorLines = $result['ErrorLines'];

$answer = strlen($answer) < 300 ? $answer : substr($answer, 0, 300) . "...";

fwrite($handle, "#  " . date('Y-m-d H:i:s') . " " . $sessionId .  " (" . $duration . "s)\n");
fwrite($handle, "Q: " . $query . "\n");
fwrite($handle, "A: " . $answer . "\n");

foreach ($optionKeys as $i => $optionKey) {
    fwrite($handle, " * " . $optionKey . " (" . $optionValues[$i] . ")\n");
}

foreach ($errorLines as $errorLine) {
    fwrite($handle, "E: " . $errorLine . "\n");
}

fwrite($handle, "\n");
fclose($handle);

header('content-type: application/json');
echo implode("\n", $output);
