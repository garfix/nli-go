<?php
/**
 * Calls "nli answer" and echoes the result in JSON.
 */

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";

session_start();
$sessionId = session_id();

$command = __DIR__ . '/nli';
$configPath = __DIR__ . '/../resources/dbpedia/config-online.json';
$fullCommand = sprintf('%s -s %s -c %s "%s"', $command, $sessionId, $configPath, $query);

$process = exec($fullCommand, $output);

$result = json_decode(implode($output), true);

$handle = fopen(__DIR__ . '/log/' . date('Y-m') . '-queries.log', 'a');

$answer = $result['Answer'];
$optionKeys = $result['OptionKeys'];
$optionValues = $result['OptionValues'];

fwrite($handle, "#  " . date('Y-m-d H:i:s') . " " . $sessionId . "\n");
fwrite($handle, "Q: " . $query . "\n");
fwrite($handle, "A: " . $answer . "\n");

foreach ($optionKeys as $i => $optionKey) {
    fwrite($handle, " * " . $optionKey . " (" . $optionValues[$i] . ")\n");
}

fwrite($handle, "\n");
fclose($handle);

header('content-type: application/json');
echo implode("\n", $output);
