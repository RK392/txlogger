import os
import datetime

# Hardcoded Target AFR string
TARGET_AFR_STRING = "20:0:147:~1:0:147:~2:0:147:~3:0:147:~4:0:147:~5:0:147:~6:0:147:~7:0:147:~8:0:147:~9:0:147:~10:0:147:~11:0:132:~12:0:130:~13:0:128:~14:0:126:~15:0:124:~16:0:122:~17:0:120:~0:1:147:~1:1:147:~2:1:147:~3:1:147:~4:1:147:~5:1:147:~6:1:147:~7:1:147:~8:1:147:~9:1:147:~10:1:147:~11:1:131:~12:1:129:~13:1:127:~14:1:125:~15:1:123:~16:1:121:~17:1:119:~0:2:147:~1:2:147:~2:2:147:~3:2:147:~4:2:147:~5:2:147:~6:2:147:~7:2:147:~8:2:147:~9:2:147:~10:2:147:~11:2:130:~12:2:128:~13:2:126:~14:2:124:~15:2:122:~16:2:120:~17:2:118:~0:3:147:~1:3:147:~2:3:147:~3:3:147:~4:3:147:~5:3:147:~6:3:147:~7:3:147:~8:3:147:~9:3:147:~10:3:147:~11:3:129:~12:3:127:~13:3:125:~14:3:123:~15:3:121:~16:3:119:~17:3:117:~0:4:147:~1:4:147:~2:4:147:~3:4:147:~4:4:147:~5:4:147:~6:4:147:~7:4:147:~8:4:147:~9:4:147:~10:4:147:~11:4:128:~12:4:126:~13:4:124:~14:4:122:~15:4:120:~16:4:118:~17:4:116:~0:5:147:~1:5:147:~2:5:147:~3:5:147:~4:5:147:~5:5:147:~6:5:147:~7:5:147:~8:5:147:~9:5:147:~10:5:147:~11:5:127:~12:5:125:~13:5:123:~14:5:121:~15:5:119:~16:5:117:~17:5:115:~0:6:147:~1:6:147:~2:6:147:~3:6:147:~4:6:147:~5:6:147:~6:6:147:~7:6:147:~8:6:147:~9:6:147:~10:6:147:~11:6:126:~12:6:124:~13:6:122:~14:6:120:~15:6:119:~16:6:117:~17:6:115:~0:7:147:~1:7:147:~2:7:147:~3:7:147:~4:7:147:~5:7:147:~6:7:147:~7:7:147:~8:7:147:~9:7:147:~10:7:147:~11:7:127:~12:7:125:~13:7:123:~14:7:121:~15:7:119:~16:7:117:~17:7:116:~0:8:147:~1:8:147:~2:8:147:~3:8:147:~4:8:147:~5:8:147:~6:8:147:~7:8:147:~8:8:147:~9:8:147:~10:8:147:~11:8:128:~12:8:126:~13:8:124:~14:8:122:~15:8:120:~16:8:119:~17:8:117:~0:9:147:~1:9:147:~2:9:147:~3:9:147:~4:9:147:~5:9:147:~6:9:147:~7:9:147:~8:9:147:~9:9:147:~10:9:147:~11:9:129:~12:9:127:~13:9:125:~14:9:123:~15:9:121:~16:9:119:~17:9:117:~0:10:147:~1:10:147:~2:10:147:~3:10:147:~4:10:147:~5:10:147:~6:10:147:~7:10:147:~8:10:147:~9:10:147:~10:10:147:~11:10:130:~12:10:128:~13:10:126:~14:10:124:~15:10:122:~16:10:120:~17:10:118:~0:11:147:~1:11:147:~2:11:147:~3:11:147:~4:11:147:~5:11:147:~6:11:147:~7:11:147:~8:11:147:~9:11:147:~10:11:147:~11:11:131:~12:11:129:~13:11:127:~14:11:125:~15:11:123:~16:11:121:~17:11:119:~0:12:147:~1:12:147:~2:12:147:~3:12:147:~4:12:147:~5:12:147:~6:12:147:~7:12:147:~8:12:147:~9:12:147:~10:12:147:~11:12:132:~12:12:130:~13:12:128:~14:12:126:~15:12:124:~16:12:122:~17:12:120:~0:13:147:~1:13:147:~2:13:147:~3:13:147:~4:13:147:~5:13:147:~6:13:147:~7:13:147:~8:13:147:~9:13:147:~10:13:147:~11:13:133:~12:13:131:~13:13:129:~14:13:127:~15:13:125:~16:13:123:~17:13:121:~0:14:147:~1:14:147:~2:14:147:~3:14:147:~4:14:147:~5:14:147:~6:14:147:~7:14:147:~8:14:147:~9:14:147:~10:14:147:~11:14:134:~12:14:132:~13:14:130:~14:14:128:~15:14:126:~16:14:124:~17:14:122:~0:15:147:~1:15:147:~2:15:147:~3:15:147:~4:15:147:~5:15:147:~6:15:147:~7:15:147:~8:15:147:~9:15:147:~10:15:147:~11:15:134:~12:15:132:~13:15:130:~14:15:129:~15:15:127:~16:15:125:~17:15:123:~"
LOG_FILENAME = "t7_fuel_log.txt"

def get_last_logged_fuel():
    if not os.path.exists(LOG_FILENAME):
        return None
    
    last_fuel = None
    with open(LOG_FILENAME, 'r') as log_file:
        for line in log_file:
            if line.startswith("NEW_FUEL:"):
                last_fuel = line.replace("NEW_FUEL:", "").strip()
    return last_fuel

def write_to_log(old_fuel, new_fuel):
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    with open(LOG_FILENAME, 'a') as log_file:
        log_file.write(f"[{timestamp}]\n")
        log_file.write(f"OLD_FUEL:{old_fuel}\n")
        log_file.write(f"NEW_FUEL:{new_fuel}\n")
        log_file.write("-" * 40 + "\n")

def parse_t7_string(data_string):
    parsed_data = {}
    
    # FIX 1: Split by the correct delimiter ':~'
    tokens = [t for t in data_string.split(':~') if t.strip()]
    
    for token in tokens:
        parts = token.split(':')
        if len(parts) == 3:
            x, y, val = parts
            x, y, val = int(x), int(y), float(val)
            
            # Normalize the '20' start coordinate to '0' so the different maps align mathematically
            if x == 20: 
                x = 0
                
            parsed_data[(x, y)] = val
            
    return parsed_data

def serialize_t7_string(data_dict):
    output = []
    
    for i, ((x, y), val) in enumerate(data_dict.items()):
        
        # FIX 2: Restore the '20' prefix for the very first element in the string
        out_x = 20 if (i == 0 and x == 0) else x
        
        if val == 0:
            val_str = "0"
        elif isinstance(val, float) and val.is_integer():
            val_str = str(int(val))
        else:
            val_str = str(val)
            
        output.append(f"{out_x}:{y}:{val_str}")
        
    return ":~".join(output) + ":~"

def calculate_new_fuel_map(input_data_str, current_fuel_str, target_afr_str, input_type=1, smoothing=0.5):
    input_data = parse_t7_string(input_data_str)
    current_fuel = parse_t7_string(current_fuel_str)
    target_afr = parse_t7_string(target_afr_str)
    
    new_fuel_map = {}
    
    for (x, y), fuel_val in current_fuel.items():
        if input_type == 2:  # Fuel Factor
            factor = input_data.get((x, y), 1024.0)
            correction_ratio = factor / 1024.0
            
            smoothed_ratio = 1.0 + (correction_ratio - 1.0) * smoothing
            new_fuel_val = round(fuel_val * smoothed_ratio, 2)
            new_fuel_map[(x, y)] = new_fuel_val
        elif input_type == 3:  # LambdaInt
            trim_val = input_data.get((x, y), 0.0)
            correction_ratio = 1.0 + (trim_val / 100.0)
            
            smoothed_ratio = 1.0 + (correction_ratio - 1.0) * smoothing
            new_fuel_val = round(fuel_val * smoothed_ratio, 2)
            new_fuel_map[(x, y)] = new_fuel_val
        else:  # Actual Lambda
            act_lam = input_data.get((x, y), 0.0)
            tgt_afr_raw = target_afr.get((x, y), 147.0)
            
            tgt_lam = tgt_afr_raw / 147.0
            
            if act_lam > 0:
                correction_ratio = act_lam / tgt_lam
                smoothed_ratio = 1.0 + (correction_ratio - 1.0) * smoothing
                new_fuel_val = round(fuel_val * smoothed_ratio, 2)
                new_fuel_map[(x, y)] = new_fuel_val
            else:
                new_fuel_map[(x, y)] = round(fuel_val, 2)
            
    return new_fuel_map

def main():
    print("--- Saab Trionic 7 Fuel Map Adjuster ---")
    
    # 1. Prompt for current fuel map
    current_fuel_str = input("Enter current fuel string (Leave blank to use last adjusted from log): ").strip()
    
    if not current_fuel_str:
        current_fuel_str = get_last_logged_fuel()
        if current_fuel_str:
            print("=> Loaded previous fuel map from log.")
        else:
            print("Error: No log found or log is empty. You must provide a current fuel string.")
            return

    # 2. Prompt for data type
    data_type = input("Do you want to use Actual Lambda (1), Fuel Factor string (2), or LambdaInt string (3)? [Default: 1]: ").strip()
    
    if data_type == '2':
        input_type = 2
        input_data_str = input("Enter fuel factor string: ").strip()
    elif data_type == '3':
        input_type = 3
        input_data_str = input("Enter lambdaint string: ").strip()
    else:
        input_type = 1
        input_data_str = input("Enter actual lambda string: ").strip()
        
    if not input_data_str:
        print("Error: Input string cannot be empty.")
        return

    # 3. Prompt for smoothing factor
    smoothing_input = input("Enter smoothing factor percentage (0-100) [Default: 50]: ").strip()
    if not smoothing_input:
        smoothing = 0.5
    else:
        try:
            smoothing = float(smoothing_input) / 100.0
        except ValueError:
            print("Invalid smoothing factor. Using default 50%.")
            smoothing = 0.5

    # 4. Calculate adjustment
    print("\nCalculating new fuel map...")
    adjusted_fuel_data = calculate_new_fuel_map(input_data_str, current_fuel_str, TARGET_AFR_STRING, input_type, smoothing)
    adjusted_fuel_string = serialize_t7_string(adjusted_fuel_data)

    # 5. Display and Log
    print("\n--- NEW FUEL MAP STRING ---")
    print(adjusted_fuel_string)
    
    write_to_log(current_fuel_str, adjusted_fuel_string)
    print(f"\n=> Map successfully written to '{LOG_FILENAME}'")

if __name__ == "__main__":
    main()