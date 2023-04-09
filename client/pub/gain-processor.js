class GainProcessor extends AudioWorkletProcessor {

  // Custom AudioParams can be defined with this static getter.
  static get parameterDescriptors() {
    return [{ name: 'gain', defaultValue: 1 }];
  }

  constructor() {
    // The super constructor call is required.
    super();
  }

  /*
  process(inputs, outputs, parameters) {
    let input = inputs[0];
    let output = outputs[0];
    let gain = parameters.gain;
    for (let channel = 0; channel < input.length; ++channel) {
      let inputChannel = input[channel];
      let outputChannel = output[channel];
      for (let i = 0; i < inputChannel.length; ++i)
        outputChannel[i] = inputChannel[i];
    }
  */
  process(inputs, outputs, parameters) {
    const input = inputs[0][0];
    const output = outputs[0][0];
    if (input == undefined) {
      return true;
    }
    const bufarr = new Int16Array(new ArrayBuffer(input.length * 2));
    for (let i = 0; i < input.length; ++i) {
      output[i] = input[i]
      bufarr[i] = input[i] < 0 ? input[i] * 0x8000 : input[i] * 0x7FFF;
    }
    if (bufarr.length > 0) {
      this.port.postMessage(bufarr);
    }
    return true;
  }
}

registerProcessor('gain-processor', GainProcessor);
