# Golang Auto Clicker

An auto clicker made in Go.

Allows you to configure multiple mouse positions for clicking, with a
main position that will be clicked more than the other positions.

Useful to clicker games.

## How to use

1. Run the program with a `-config` flag to configure the positions to
   click.

    ```bash
    $ cd clicker
    $ go run . -config
    ```

2. Click at the positions that the clicker will click. The first click
   indicates the main position, that will be the most clicked.
3. When all positions are configured, right click to end the
   configuration.
4. Run the program.
   ```bash
   $ go run .
   ```
   Press any key to stop the clicker.

The clicker clicks 15 times at the main position, then clicks once
at the second, back to the main and clicks more 15 times, go to the
third, and so on...
