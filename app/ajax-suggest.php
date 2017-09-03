<?php
/**
 * Calls "nli suggest" and echoes the result in JSON.
 */

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";
$tokens = explode(',', $query);
$sentence = implode(' ', $tokens);

$command = __DIR__ . '/nli suggest';
$configPath = __DIR__ . '/../resources/dbpedia/config-online.json';
$fullCommand = sprintf('%s %s "%s"', $command, $configPath, $sentence);

$process = exec($fullCommand, $output);

header('content-type: application/json');
echo implode("\n", $output);
