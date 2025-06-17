Bugfix: Add missing value for consumption metrics

Until now the value for the consumption metrics was missing any assignment of a
real value as the variable have only been defined with the default zero value.
Beside that you are also able to add a currency label to make sure you can
calculate correct currencies.

https://github.com/promhippie/scw_exporter/issues/138
