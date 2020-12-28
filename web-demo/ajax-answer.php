<?php

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";
$app = isset($_REQUEST['app']) ? $_REQUEST['app'] : "dbpedia";

if (!in_array($app, ['dbpedia', 'blocks'])) {
    die('Unsupported app');
}
session_start();
$sessionId = session_id();

$command = __DIR__ . '/../bin/nli';
$configPath = __DIR__ . '/../resources/' . $app;
$varDir = __DIR__ . '/../var';

// query

$start = microtime(true);
$fullCommand = sprintf('%s -s %s -c %s -d %s -r json "%s"', $command, $sessionId, $configPath, $varDir, $query);
exec($fullCommand, $output);
$end = microtime(true);
$duration = sprintf("%.2f", $end - $start);

$result = json_decode(implode($output), true);
$response = json_encode($result, JSON_PRETTY_PRINT);

logResult($result, $duration, $sessionId);

header('content-type: application/json');
echo $response;


function logResult($result, $duration, $sessionId)
{
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
}